package config

import "path/filepath"

const (
	uploadsDir      = "/uploads"
	photosDir       = "/photos"
	swaggerFilePath = "/api/swagger.yml"
)

type PathConfig struct {
	RootDir         string
	UploadDirs      *UploadDirs
	SwaggerFilePath string
}

type UploadDirs struct {
	UploadDir string
	PhotoDir  string
}

func (pc *PathConfig) FullPath(paths ...string) string {
	path := filepath.Join(paths...)

	return filepath.Join(pc.RootDir, path)
}

func pathFromEnv(rootDir string) *PathConfig {
	return &PathConfig{
		rootDir,
		&UploadDirs{uploadsDir, photosDir},
		swaggerFilePath,
	}
}
