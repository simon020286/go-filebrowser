package acutime

import (
	"os"
	"syscall"
	"time"
)

// Ctime creation time.
func ctime(fi os.FileInfo) time.Time {
	return time.Unix(0, fi.Sys().(*syscall.Win32FileAttributeData).CreationTime.Nanoseconds())
}

// Utime last write time.
func utime(fi os.FileInfo) time.Time {
	return time.Unix(0, fi.Sys().(*syscall.Win32FileAttributeData).LastWriteTime.Nanoseconds())
}

// Atime last access time.
func atime(fi os.FileInfo) time.Time {
	return time.Unix(0, fi.Sys().(*syscall.Win32FileAttributeData).LastAccessTime.Nanoseconds())
}
