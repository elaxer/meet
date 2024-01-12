package router

import (
	"meet/internal/config"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type fileSystem struct {
	pathConfig *config.PathConfig
	fs         http.FileSystem
}

type file struct {
	http.File
}

func (f file) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (fs *fileSystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	return &file{f}, nil
}

func (c *configurator) configureFileServer(r *mux.Router) {
	fs := &fileSystem{c.pathConfig, http.Dir(c.pathConfig.FullPath(c.pathConfig.UploadDirs.UploadDir))}

	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(fs)))
}
