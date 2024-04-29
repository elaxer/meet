package helper

import (
	"meet/internal/config"
	"meet/internal/pkg/app"
	"path/filepath"
	"strconv"
)

type PathHelper interface {
	UploadPath(filename string, uploadType app.UploadType, userID int) string
}

type pathHelper struct {
	pathConfig *config.PathConfig
}

func NewPathHelper(pathConfig *config.PathConfig) PathHelper {
	return &pathHelper{pathConfig}
}

func (ph *pathHelper) UploadPath(filename string, uploadType app.UploadType, userID int) string {
	return filepath.Join(
		ph.pathConfig.RootDir,
		ph.pathConfig.UploadDirs.UploadDir,
		string(uploadType),
		strconv.Itoa(userID),
		filename,
	)
}
