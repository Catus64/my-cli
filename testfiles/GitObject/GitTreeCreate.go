package GitObjLib

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	gitpath "gocmd/testfiles/Gitrepostruct"
	logger "gocmd/testfiles/Helper"
	"os"
	"path/filepath"
	"sort"
)

type treeEntry struct {
	Mode   string
	Name   string
	SHA    string
	isTree bool
}

func TreeFromIndex(repo gitpath.GitRepository, index GitIndex) (string, error) {
	// Phase 1: group entries by their parent directory
	// key = directory path, value = list of entries in that dir
	contents := make(map[string][]treeEntry)
	contents[""] = []treeEntry{} // ensure root always exists

	for _, entry := range index.Entries {
		dir := filepath.ToSlash(filepath.Dir(entry.Name))
		if dir == "." {
			dir = ""
		}

		// add the entry to its direct parent
		mode := fmt.Sprintf("%02o%04o", entry.ModeType, entry.ModePerms)
		logger.L().Info("Adding index entry", "Name", entry.Name, "Dir", dir, "Mode", mode)
		contents[dir] = append(contents[dir], treeEntry{
			Mode: mode,
			Name: filepath.Base(entry.Name),
			SHA:  entry.SHA,
		})

		// ensure ALL ancestor directories exist in the map
		// so parent trees get created even if they hold no direct files
		key := dir
		for key != "" {
			key = filepath.ToSlash(filepath.Dir(key))
			if key == "." {
				key = ""
			}
			if _, exists := contents[key]; !exists {
				fmt.Printf("Ensuring ancestor directory %q exists in tree map\n", key)
				contents[key] = []treeEntry{}
			}
		}
	}

	// Phase 2: collect and sort paths by length descending
	// children always come before parents this way
	paths := make([]string, 0, len(contents))
	for path := range contents {
		paths = append(paths, path)
	}
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) > len(paths[j])
	})

	for path := range paths {
		logger.L().Debug("treepath", "Tree path in order of creation:", paths[path])
	}

	// Phase 3: build trees bottom-up
	var rootSHA string
	for _, path := range paths {
		entries := contents[path]

		// sort entries within this tree
		// git requires entries sorted by name, but dirs sorted as if they have trailing /
		sort.Slice(entries, func(i, j int) bool {
			nameI := entries[i].Name
			nameJ := entries[j].Name
			if entries[i].isTree {
				nameI += "/"
			}
			if entries[j].isTree {
				nameJ += "/"
			}
			return nameI < nameJ
		})

		// build the raw tree binary content
		sha, err := writeTreeObject(repo, entries)
		if err != nil {
			return "", fmt.Errorf("failed to write tree for %q: %w", path, err)
		}

		rootSHA = sha

		// attach this tree to its parent
		if path != "" {
			parent := filepath.ToSlash(filepath.Dir(path))
			if parent == "." {
				parent = ""
			}
			base := filepath.Base(path)
			contents[parent] = append(contents[parent], treeEntry{
				Mode:   "40000", // directory mode — no leading zero
				Name:   base,
				SHA:    sha,
				isTree: true,
			})
		}
	}

	return rootSHA, nil
}

func writeTreeObject(repo gitpath.GitRepository, entries []treeEntry) (string, error) {
	var body bytes.Buffer

	for _, e := range entries {
		// format: "<mode> <name>\0<20-byte-raw-sha>"
		body.WriteString(e.Mode)
		body.WriteByte(' ')
		body.WriteString(e.Name)
		body.WriteByte(0x00)

		// SHA must be raw 20 bytes, not hex string
		shaBytes, err := hex.DecodeString(e.SHA)
		if err != nil {
			return "", fmt.Errorf("invalid sha %q: %w", e.SHA, err)
		}
		body.Write(shaBytes)
	}

	// git object format: "<type> <size>\0<content>"
	var object bytes.Buffer
	object.WriteString(fmt.Sprintf("tree %d\x00", body.Len()))
	object.Write(body.Bytes())

	// SHA1 of the full object
	raw := object.Bytes()
	hash := sha1.Sum(raw)
	hexSHA := hex.EncodeToString(hash[:])

	// custom write function temporary here to make it work
	if err := writeObject(repo, hexSHA, raw); err != nil {
		return "", err
	}

	return hexSHA, nil
}

// writeObject zlib-compresses and writes to the object store
func writeObject(repo gitpath.GitRepository, hexSHA string, raw []byte) error {
	location_1 := hexSHA[:2]
	location_2 := hexSHA[2:]
	path := gitpath.Repo_Path(repo, "objects", location_1, location_2)

	logger.L().Debug("Writing tree object", hexSHA, path)

	// don't rewrite if it already exists
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := zlib.NewWriter(f)
	defer w.Close()

	_, err = w.Write(raw)
	return err
}
