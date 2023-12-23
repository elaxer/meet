package controller

import (
	"meet/internal/app"
	"meet/internal/app/server"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
)

type swaggerController struct {
	cfg *app.Config
}

func newSwaggerController(cfg *app.Config) *swaggerController {
	return &swaggerController{cfg}
}

func (sc *swaggerController) Api(w http.ResponseWriter, r *http.Request) {
	b, err := os.ReadFile(filepath.Join(app.RootDir, sc.cfg.SwaggerFile))
	if err != nil {
		server.ResponseError(w, err, http.StatusInternalServerError)

		return
	}

	swaggerJSON, err := yaml.YAMLToJSON(b)
	if err != nil {
		server.ResponseError(w, err, http.StatusInternalServerError)

		return
	}

	server.ResponseRaw(w, swaggerJSON, http.StatusOK)
}
