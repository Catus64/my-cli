package GitObjLib

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"os"
)

type GitIndexEntry struct {
	CtimeSec    [2]uint32 // seconds part of ctime
	CtimeNano   int32     // nanoseconds part of ctime
	MtimeSec    [2]uint32 // seconds part of mtime
	MtimeNano   int32     // nanoseconds part of mtime
	Dev         uint32    // device ID
	Ino         uint32    // inode
	ModeType    uint16    // Git object type bits: regular/symlink/gitlink
	ModePerms   uint16    // File permissions (0644/0755)
	UID         uint32    // User ID
	GID         uint32    // Group ID
	FSize       uint32    // File size in bytes
	SHA         string    // SHA-1 hash of the file content
	AssumeValid bool      // skip stat check
	Stage       uint16    // stage number: 0-3
	Name        string    // file path relative to repo root
}

type GitIndex struct {
	Version int
	Entries []GitIndexEntry
}

func NewGitIndex() *GitIndex {
	return &GitIndex{
		Version: 2,
		Entries: []GitIndexEntry{},
	}
}

func NewGitIndexWithEntries(entries []GitIndexEntry) *GitIndex {
	return &GitIndex{
		Version: 2,
		Entries: entries,
	}
}

func Index_Read(repo gitpath.GitRepository) (*GitIndex, error) {
	// Find index file path
	indexPath := gitpath.Repo_Path(repo, "index")

	// Check if index file exists
	// create new index with empty entries if it doesn't exist
	if _, err := os.Stat(indexPath); err != nil {
		if os.IsNotExist(err) {
			index := NewGitIndex()
			return index, nil
		}
		return nil, fmt.Errorf("error accessing save list: %v", err)
	}

	// Read index file content
	raw, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("error reading save list: %v", err)
	}

	header := raw[:12]

	//sig should be "DIRC"
	signature := header[:4]
	if string(signature) != "DIRC" {
		return nil, fmt.Errorf("invalid index file: wrong signature")
	}

	//version should be 2
	//get last byte of version to check if it's 2
	//only support ver 2 for now
	version := header[4:8]
	version_number := version[len(version)-1]
	if int(version_number) != 2 {
		return nil, fmt.Errorf("unsupported index version: %s", version)
	}

	//get file count (last byte of header)
	file_count := header[8:12]
	count := int(file_count[len(file_count)-1])

	entries := []GitIndexEntry{}

	fmt.Println("File count: ", count)
	fmt.Println("Signature: ", string(signature))
	fmt.Println("Version: ", version_number)

	content := raw[12:]
	idx := 0

	//read entry
	for i := 0; i < count; i++ {

		ctimeS := binary.BigEndian.Uint32(content[idx : idx+4])
		ctimeNs := binary.BigEndian.Uint32(content[idx+4 : idx+8])

		mtimeS := binary.BigEndian.Uint32(content[idx+8 : idx+12])
		mtimeNs := binary.BigEndian.Uint32(content[idx+12 : idx+16])

		dev := binary.BigEndian.Uint32(content[idx+16 : idx+20])
		ino := binary.BigEndian.Uint32(content[idx+20 : idx+24])

		// unused field (should always be 0)
		unused := binary.BigEndian.Uint16(content[idx+24 : idx+26])
		if unused != 0 {
			return nil, fmt.Errorf("unexpected value in unused field")
		}

		mode := binary.BigEndian.Uint16(content[idx+26 : idx+28])

		// upper 4 bits describe object type
		modeType := mode >> 12

		if modeType != 0b1000 &&
			modeType != 0b1010 &&
			modeType != 0b1110 {
			return nil, fmt.Errorf("invalid mode type")
		}

		// lower 9 bits are permissions
		modePerms := mode & 0b0000000111111111

		uid := binary.BigEndian.Uint32(content[idx+28 : idx+32])
		gid := binary.BigEndian.Uint32(content[idx+32 : idx+36])

		fsize := binary.BigEndian.Uint32(content[idx+36 : idx+40])

		shaBytes := content[idx+40 : idx+60]
		sha := hex.EncodeToString(shaBytes)

		flags := binary.BigEndian.Uint16(content[idx+60 : idx+62])

		// Highest bit
		flagAssumeValid := (flags & 0b1000000000000000) != 0

		// Extended flag
		flagExtended := (flags & 0b0100000000000000) != 0
		if flagExtended {
			return nil, fmt.Errorf("extended flags not supported")
		}

		// Merge stage (2 bits)
		flagStage := flags & 0b0011000000000000

		// File name length (lower 12 bits)
		nameLength := flags & 0b0000111111111111

		idx += 62

		var rawName []byte

		if nameLength < 0xFFF {

			if content[idx+int(nameLength)] != 0 {
				return nil, fmt.Errorf("name not null terminated")
			}

			rawName = content[idx : idx+int(nameLength)]

			idx += int(nameLength) + 1

		} else {

			// Very long path case
			nullIdx := bytes.IndexByte(content[idx+0xFFF:], 0)
			if nullIdx == -1 {
				return nil, fmt.Errorf("could not find null terminator")
			}

			nullIdx += idx + 0xFFF

			rawName = content[idx:nullIdx]

			idx = nullIdx + 1
		}

		name := string(rawName)

		for idx%8 != 0 {
			idx++
		}
		entry := GitIndexEntry{
			CtimeSec: [2]uint32{ctimeS, ctimeNs},
			MtimeSec: [2]uint32{mtimeS, mtimeNs},
			Dev:      dev,
			Ino:      ino,

			ModeType:  modeType,
			ModePerms: modePerms,

			UID: uid,
			GID: gid,

			FSize: fsize,
			SHA:   sha,

			AssumeValid: flagAssumeValid,
			Stage:       flagStage,

			Name: name,
		}

		print(entry.Name, " ", entry.SHA, " ", entry.ModeType, " ", entry.ModePerms, "\n")

		entries = append(entries, entry)

	}

	return nil, nil
}
