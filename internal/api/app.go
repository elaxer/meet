package api

import (
	"fmt"
	"meet/internal/api/middleware"
	"meet/internal/api/router"
	"meet/internal/pkg/app"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	config *app.Config
}

func NewApp() *App {
	return &App{
		app.NewConfig(),
	}
}

func (a *App) Run() {
	r := mux.NewRouter()
	router.Configure(r, controllers, middlewares)

	http.Handle(
		"/",
		middleware.LogMiddleware(
			middleware.CORSMiddleware(
				middleware.ContentLength(
					middleware.FileSize(r),
				),
			),
		),
	)

	addr := fmt.Sprintf("%s:%s", a.config.ServerConfig.Host, a.config.ServerConfig.Port)

	http.ListenAndServe(addr, nil)
}
