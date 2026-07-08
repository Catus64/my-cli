package GitObjLib

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"os"
)

func Index_Write(repo gitpath.GitRepository, index GitIndex) error {

	buf := new(bytes.Buffer)

	// Git writes DIRC as a header
	//if no "DIRC" git consider index invalid
	buf.WriteString("DIRC")
	binary.Write(buf, binary.BigEndian, index.Version)
	binary.Write(buf, binary.BigEndian, uint32(len(index.Entries)))

	// Write each entry
	for _, e := range index.Entries {
		//entryStart := buf.Len()

		// Timestamp related entry fields (Ctime and Mtime)
		binary.Write(buf, binary.BigEndian, e.CtimeSec)
		binary.Write(buf, binary.BigEndian, e.CtimeNano)
		binary.Write(buf, binary.BigEndian, e.MtimeSec)
		binary.Write(buf, binary.BigEndian, e.MtimeNano)

		// Write Dev(device) and Ino(inode) numbers
		// refer: GitIndexEntry for any of the field details here
		binary.Write(buf, binary.BigEndian, e.Dev)
		binary.Write(buf, binary.BigEndian, e.Ino)

		// Mode Wrute ModeType and ModePerms into a 32-bit field
		mode := (uint32(e.ModeType) << 12) | uint32(e.ModePerms)
		binary.Write(buf, binary.BigEndian, mode)

		binary.Write(buf, binary.BigEndian, e.UID)
		binary.Write(buf, binary.BigEndian, e.GID)
		binary.Write(buf, binary.BigEndian, e.FSize)

		// hex string back to 20 raw bytes from the 40 character hex string
		shaBytes, _ := hex.DecodeString(e.SHA)
		buf.Write(shaBytes)

		// Flags contains AssumeValid + Stage + Name(length) into 2 bytes
		nameBytes := []byte(e.Name)
		nameLen := len(nameBytes)
		if nameLen >= 0xFFF {
			nameLen = 0xFFF
		}
		var assumeValid uint16
		if e.AssumeValid {
			assumeValid = 0x1 << 15
		}
		flags := assumeValid | e.Stage | uint16(nameLen)
		binary.Write(buf, binary.BigEndian, flags)

		// Name + null terminator
		buf.Write(nameBytes)
		buf.WriteByte(0x00)

		// Padding to 8-byte boundary
		// 62 = all fixed fields, +len(name) +1 null byte
		//
		entryLen := 62 + len(nameBytes) + 1
		if entryLen%8 != 0 {
			pad := 8 - (entryLen % 8)
			buf.Write(make([]byte, pad))
		}

		// nameBytes = []byte(e.Name)
		// entryLen = 62 + len(nameBytes) + 1

		// pad := 0
		// if entryLen%8 != 0 {
		// 	pad = 8 - (entryLen % 8)
		// }

		// totalEntry := entryLen + pad

		// fmt.Printf("entry=%-20s  fixedPlusName=%d  pad=%d  total=%d  bufLen=%d\n",
		// 	e.Name, entryLen, pad, totalEntry, buf.Len()-entryStart)
	}

	checksum := sha1.Sum(buf.Bytes())
	buf.Write(checksum[:])

	// Write buffer to file in one shot
	path := gitpath.Repo_Path(repo, "index")
	return os.WriteFile(path, buf.Bytes(), 0644)
}
