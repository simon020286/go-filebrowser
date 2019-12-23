package utils

import (
	"filebrowser/acutime"
	"filebrowser/models"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/chebyrash/promise"
)

// IFileUtils interface.
type IFileUtils interface {
	GetFiles(root string) ([]models.File, error)
	CreateFolder(src string) error
	Copy(src, dst string) (int64, error)
	CopyPromise(src, dst string) *promise.Promise
	Delete(src string) error
	DeletePromise(src string) *promise.Promise
	Move(src, dst string) (int64, error)
	GetAbsolutePath(relative string) string
	Save(file io.Reader, name string, path string) (int64, error)
	SavePromise(file io.Reader, name string, path string) *promise.Promise
}

type fileUtils struct {
	basePath string
}

func (fu *fileUtils) GetAbsolutePath(relative string) string {
	return path.Join(fu.basePath, relative)
}

// GetFiles list.
func (fu fileUtils) GetFiles(root string) ([]models.File, error) {
	root = fu.GetAbsolutePath(root)
	filesInfo, err := ioutil.ReadDir(root)
	files := make([]models.File, 0)

	if err != nil {
		return nil, err
	}

	for _, fileInfo := range filesInfo {
		files = append(files, models.File{Name: fileInfo.Name(), Size: fileInfo.Size(), IsDir: fileInfo.IsDir(), CreatedAt: acutime.Ctime(fileInfo), UpdatedAt: acutime.Utime(fileInfo)})
	}

	return files, nil
}

// CreateFolder in specific path.
func (fu fileUtils) CreateFolder(src string) error {
	src = fu.GetAbsolutePath(src)
	return os.Mkdir(src, os.ModePerm)
}

func (fu fileUtils) Move(src, dst string) (int64, error) {
	src = fu.GetAbsolutePath(src)
	dst = fu.GetAbsolutePath(dst)
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}
	err = os.Rename(src, dst)
	if err != nil {
		return 0, err
	}

	return sourceFileStat.Size(), nil
}

// Copy file or folder.
func (fu fileUtils) Copy(src, dst string) (int64, error) {
	src = fu.GetAbsolutePath(src)
	dst = fu.GetAbsolutePath(dst)
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func (fu fileUtils) CopyPromise(src, dst string) *promise.Promise {
	return promise.New(func(resolve func(interface{}), reject func(error)) {
		size, err := fu.Copy(src, dst)
		if err != nil {
			reject(err)
			return
		}

		resolve(size)
	})
}

// Delete file or folder.
func (fu fileUtils) Delete(src string) error {
	src = fu.GetAbsolutePath(src)
	_, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.RemoveAll(src)
}

func (fu fileUtils) DeletePromise(src string) *promise.Promise {
	return promise.New(func(resolve func(interface{}), reject func(error)) {
		err := fu.Delete(src)
		if err != nil {
			reject(err)
			return
		}
		resolve(nil)
	})
}

func (fu fileUtils) Save(file io.Reader, name string, dst string) (int64, error) {
	dst = path.Join(fu.GetAbsolutePath(dst), name)

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	return io.Copy(out, file)
}

func (fu fileUtils) SavePromise(file io.Reader, name string, dst string) *promise.Promise {
	return promise.New(func(resolve func(interface{}), reject func(error)) {
		size, err := fu.Save(file, name, dst)
		if err != nil {
			reject(err)
			return
		}

		resolve(size)
	})
}

// NewFileUtils constructor.
func NewFileUtils(basePath string) IFileUtils {
	return &fileUtils{basePath: basePath}
}
