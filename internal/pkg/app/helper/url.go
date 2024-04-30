package helper

import (
	"fmt"
	"meet/internal/config"
	"meet/internal/pkg/app"
	"net/url"
	"strings"
)

type URLHelper interface {
	UploadURL(filename string, uploadType app.UploadType, userID int) string
}

type urlHelper struct {
	serverConfig *config.ServerConfig
	uploadDirs   *config.UploadDirs
}

func NewURLHelper(serverConfig *config.ServerConfig, uploadDirs *config.UploadDirs) URLHelper {
	return &urlHelper{serverConfig, uploadDirs}
}

func (ph *urlHelper) UploadURL(filename string, uploadType app.UploadType, userID int) string {
	schemeAndHost := strings.Split(ph.serverConfig.Host, "://")
	scheme, host := schemeAndHost[0], schemeAndHost[1]

	if port := ph.serverConfig.Port; port != 80 {
		host = fmt.Sprintf("%s:%d", host, port)
	}

	u := &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   fmt.Sprintf("%s/%s/%d/%s", ph.uploadDirs.UploadDir, uploadType, userID, filename),
	}

	return u.String()
}
