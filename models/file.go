package models

import (
	"encoding/json"
	"time"
)

// File struct.
type File struct {
	Name      string
	Size      int64
	IsDir     bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MarshalJSON function.
func (file *File) MarshalJSON() (text []byte, err error) {
	fileType := "file"
	if file.IsDir {
		fileType = "dir"
	}
	s := struct {
		Name      string    `json:"name"`
		Size      int64     `json:"size"`
		FileType  string    `json:"type"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}{
		Name:      file.Name,
		Size:      file.Size,
		FileType:  fileType,
		CreatedAt: file.CreatedAt,
		UpdatedAt: file.UpdatedAt,
	}

	return json.Marshal(s)
}
