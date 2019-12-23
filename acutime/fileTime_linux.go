package acutime

import (
	"os"
	"syscall"
	"time"
)

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

// Ctime creation time.
func ctime(fi os.FileInfo) time.Time {
	return timespecToTime(fi.Sys().(*syscall.Stat_t).Ctim)
}

// Utime last write time.
func utime(fi os.FileInfo) time.Time {
	return timespecToTime(fi.Sys().(*syscall.Stat_t).Mtim)
}

// Atime last access time.
func atime(fi os.FileInfo) time.Time {
	return timespecToTime(fi.Sys().(*syscall.Stat_t).Atim)
}
