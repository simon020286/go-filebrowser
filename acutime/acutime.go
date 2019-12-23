package acutime

import (
	"os"
	"time"
)

// Atime last access time.
func Atime(fi os.FileInfo) time.Time {
	return atime(fi)
}

// Ctime creation time.
func Ctime(fi os.FileInfo) time.Time {
	return ctime(fi)
}

// Utime last write time.
func Utime(fi os.FileInfo) time.Time {
	return utime(fi)
}
