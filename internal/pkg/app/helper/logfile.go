package helper

import (
	"io"
	"meet/internal/pkg/app"
	"os"
)

func OpenLogFile(rootDir string) (io.WriteCloser, error) {
	return os.OpenFile(rootDir+"/"+app.LogFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
}
