//go:build windows

package gitsave

import (
	"os"

	gitobj "gocmd/testfiles/GitObject"
)

func updateStatTimes(index *gitobj.GitIndex, i int, info os.FileInfo) {
	// Windows does not expose ctime via syscall.Stat_t
	// use mtime as fallback
	index.Entries[i].CtimeSec = uint32(info.ModTime().Unix())
	index.Entries[i].CtimeNano = uint32(info.ModTime().Nanosecond())
}
