package slogger

import (
	"io"
	"os"
)

const LogFilename = "logs/log.log"

func OpenLog(rootDir string) (io.WriteCloser, error) {
	return os.OpenFile(rootDir+"/"+LogFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
}
func MustOpenLog(rootDir string) io.WriteCloser {
	f, err := OpenLog(rootDir)
	if err != nil {
		panic(err)
	}

	return f
}
