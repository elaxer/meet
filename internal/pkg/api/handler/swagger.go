package handler

import (
	"meet/internal/config"
	"meet/internal/pkg/api"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
)

type SwaggerHandler interface {
	Api(w http.ResponseWriter, r *http.Request)
}

type swaggerHandler struct {
	pathsConfig *config.PathConfig
}

func NewSwaggerHandler(pathsConfig *config.PathConfig) SwaggerHandler {
	return &swaggerHandler{pathsConfig}
}

func (sh *swaggerHandler) Api(w http.ResponseWriter, r *http.Request) {
	b, err := os.ReadFile(filepath.Join(sh.pathsConfig.RootDir, sh.pathsConfig.SwaggerFilePath))
	if err != nil {
		api.ResponseError(w, err, http.StatusInternalServerError)

		return
	}

	swaggerJSON, err := yaml.YAMLToJSON(b)
	if err != nil {
		api.ResponseError(w, err, http.StatusInternalServerError)

		return
	}

	api.ResponseRaw(w, swaggerJSON, http.StatusOK)
}
