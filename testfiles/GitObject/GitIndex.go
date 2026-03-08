package GitObjLib

type GitIndexEntry struct {
	CtimeSec    int64    // seconds part of ctime
	CtimeNano   int32    // nanoseconds part of ctime
	MtimeSec    int64    // seconds part of mtime
	MtimeNano   int32    // nanoseconds part of mtime
	Dev         uint32   // device ID
	Ino         uint32   // inode
	ModeType    uint16   // Git object type bits: regular/symlink/gitlink
	ModePerms   uint16   // File permissions (0644/0755)
	UID         uint32   // User ID
	GID         uint32   // Group ID
	FSize       uint32   // File size in bytes
	SHA         [20]byte // SHA-1 hash of the file content
	AssumeValid bool     // skip stat check
	Stage       uint8    // stage number: 0-3
	Name        string   // file path relative to repo root
}
