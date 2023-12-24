package service

import (
	"errors"
	"io"
	"meet/internal/pkg/app"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/h2non/filetype"
)

var (
	ErrFileTypeWrong = errors.New("неверный тип выгружаемого файла")
)

type FileService struct {
	uploadDir string
}

func newFileService(uploadDir string) *FileService {
	return &FileService{uploadDir}
}

func (fs *FileService) Upload(file io.ReadSeeker, types []string, uploadsSubDir string) (*os.File, string, error) {
	t, err := filetype.MatchReader(file)
	if err != nil {
		return nil, "", err
	}

	if !fs.checkType(t.MIME.Type, types) {
		return nil, "", ErrFileTypeWrong
	}

	fname := fs.generateName(t.Extension)
	path := fs.FullPath(fname, uploadsSubDir)

	f, err := os.Create(path)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	file.Seek(0, io.SeekStart)
	_, err = io.Copy(f, file)
	if err != nil {
		return nil, "", err
	}

	return f, fname, nil
}

func (fs *FileService) FullPath(fname string, uploadsSubDir string) string {
	return filepath.Join(app.RootDir, fs.uploadDir, uploadsSubDir, fname)
}

func (fs *FileService) checkType(t string, types []string) bool {
	if len(types) == 0 {
		return true
	}

	for _, v := range types {
		if v == t {
			return true
		}
	}

	return false
}

func (fs *FileService) generateName(ext string) string {
	return uuid.New().String() + "." + ext
}
