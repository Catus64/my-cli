//go:build !windows

package gitsave

import (
	"os"
	"syscall"

	gitobj "gocmd/testfiles/GitObject"
)

func updateStatTimes(index *gitobj.GitIndex, i int, info os.FileInfo) {
	if sysStat, ok := info.Sys().(*syscall.Stat_t); ok {
		index.Entries[i].CtimeSec = uint32(sysStat.Ctim.Sec)
		index.Entries[i].CtimeNano = uint32(sysStat.Ctim.Nsec)
	}
}
