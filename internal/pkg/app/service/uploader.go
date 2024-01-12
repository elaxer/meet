package service

import (
	"errors"
	"io"
	"meet/internal/config"
	"meet/internal/pkg/app/helper"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/h2non/filetype"
)

var (
	ErrFileTypeWrong = errors.New("неверный тип выгружаемого файла")
)

type FileUploaderService interface {
	Upload(file io.ReadSeeker, uploadType config.UploadType, userID int) (*os.File, error)
}

type fileUploaderService struct {
	pathHelper  helper.PathHelper
	pathsConfig *config.PathConfig
}

func NewFileUploaderService(pathHelper helper.PathHelper, pathsConfig *config.PathConfig) FileUploaderService {
	return &fileUploaderService{pathHelper, pathsConfig}
}

func (fus *fileUploaderService) Upload(file io.ReadSeeker, uploadType config.UploadType, userID int) (*os.File, error) {
	t, err := filetype.MatchReader(file)
	if err != nil {
		return nil, err
	}
	if t.MIME.Type != string(uploadType) {
		return nil, ErrFileTypeWrong
	}

	filename := fus.generateName(t.Extension)
	path := fus.pathHelper.UploadPath(filename, uploadType, userID)

	dir, _ := filepath.Split(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return nil, err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	file.Seek(0, io.SeekStart)
	if _, err := io.Copy(f, file); err != nil {
		return nil, err
	}

	return f, nil
}

func (fs *fileUploaderService) generateName(ext string) string {
	return uuid.New().String() + "." + ext
}
