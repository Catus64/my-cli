package GitObjLib

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	gitpath "gocmd/testfiles/Gitrepostruct"
	logger "gocmd/testfiles/Helper"
	"io"
	"os"
)

type GitIndexEntry struct {
	CtimeSec    uint32 // seconds part of ctime (when metadata was last changed)
	CtimeNano   uint32 // nanoseconds part of ctime
	MtimeSec    uint32 // seconds part of mtime (when content was last modified)
	MtimeNano   uint32 // nanoseconds part of mtime
	Dev         uint32 // device ID
	Ino         uint32 // inode
	ModeType    uint16 // Git object type bits: regular/symlink/gitlink
	ModePerms   uint16 // File permissions (0644/0755)
	UID         uint32 // User ID
	GID         uint32 // Group ID
	FSize       uint32 // File size in bytes
	SHA         string // SHA-1 hash of the file content
	AssumeValid bool   // skip stat check
	Stage       uint16 // stage number: 0-3
	Name        string // file path relative to repo root
}

type indexEntryHeader struct {
	CtimeS   uint32
	CtimeNs  uint32
	MtimeS   uint32
	MtimeNs  uint32
	Dev      uint32
	Ino      uint32
	Unused   uint16
	Mode     uint16
	UID      uint32
	GID      uint32
	FileSize uint32
	SHA      [20]byte
	Flags    uint16
}

type GitIndex struct {
	Version uint32
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

func Index_Read2(repo gitpath.GitRepository) (*GitIndex, error) {
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

	file, err := os.Open(indexPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	//Resolving Header

	var signature [4]byte
	err = binary.Read(reader, binary.BigEndian, &signature)
	if err != nil {
		return nil, err
	}

	if string(signature[:]) != "DIRC" {
		return nil, fmt.Errorf("invalid")
	}

	var version uint32
	err = binary.Read(reader, binary.BigEndian, &version)
	if err != nil {
		return nil, err
	}

	if version != 2 {
		return nil, fmt.Errorf("unsupported savelist version")
	}

	var count uint32
	err = binary.Read(reader, binary.BigEndian, &count)
	if err != nil {
		return nil, err
	}

	logger.L().Debug("Index header read",
		"signature", string(signature[:]),
		"version", version,
		"entry_count", count,
	)

	entries := []GitIndexEntry{}

	//Reading each entry

	for i := 0; i < int(count); i++ {
		var header indexEntryHeader
		err = binary.Read(reader, binary.BigEndian, &header)
		if err != nil {
			return nil, err
		}

		if header.Unused != 0 {
			return nil, fmt.Errorf("unexpected value in unused field")
		}

		modeType := header.Mode >> 12
		modePerms := header.Mode & 0x01FF

		flags := header.Flags

		flagAssumeValid := (flags & 0x8000) != 0
		flagExtended := (flags & 0x4000) != 0
		if flagExtended {
			return nil, fmt.Errorf("extended flags not supported")
		}

		flagStage := (flags >> 12) & 0x3000
		nameLength := flags & 0x0FFF

		// Read file name

		var nameBytes []byte

		if nameLength < 0xFFF {

			nameBytes = make([]byte, nameLength)

			_, err = io.ReadFull(reader, nameBytes)
			if err != nil {
				return nil, err
			}

			_, err = reader.ReadByte()
			if err != nil {
				return nil, err
			}

		} else {

			for {
				b, err := reader.ReadByte()
				if err != nil {
					return nil, err
				}
				if b == 0 {
					break
				}
				nameBytes = append(nameBytes, b)
			}

		}

		name := string(nameBytes)

		// Align to 8 bytes

		entrySize := 62 + len(nameBytes) + 1
		padding := (8 - (entrySize % 8)) % 8

		if padding > 0 {
			_, err = io.CopyN(io.Discard, reader, int64(padding))
			if err != nil {
				return nil, err
			}
		}

		sha := hex.EncodeToString(header.SHA[:])

		entry := GitIndexEntry{
			CtimeSec:  header.CtimeS,
			CtimeNano: header.CtimeNs,
			MtimeSec:  header.MtimeS,
			MtimeNano: header.MtimeNs,
			Dev:       header.Dev,
			Ino:       header.Ino,

			ModeType:  modeType,
			ModePerms: modePerms,

			UID: header.UID,
			GID: header.GID,

			FSize: header.FileSize,
			SHA:   sha,

			AssumeValid: flagAssumeValid,
			Stage:       flagStage,

			Name: name,
		}
		entries = append(entries, entry)
	}
	return &GitIndex{
		Version: version,
		Entries: entries,
	}, nil
}
