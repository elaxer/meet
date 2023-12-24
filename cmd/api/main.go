package main

import (
	"database/sql"
	"fmt"
	"meet/internal/api/controller"
	"meet/internal/api/middleware"
	"meet/internal/api/router"
	"meet/internal/pkg/app"
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/app/service"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var cfg *app.Config

var (
	repositories *repository.RepositoryContainer
	services     *service.ServiceContainer
	controllers  *controller.ControllerContainer
	middlewares  *middleware.MiddlewareContainer
)

func init() {
	err := godotenv.Load(app.RootDir+"/.env.local", app.RootDir+"/.env")
	if err != nil {
		panic(err)
	}

	cfg = app.NewConfig()

	db := newDB()

	repositories = repository.NewRepositoryContainer(db)
	services = service.NewServiceContainer(cfg, repositories)
	controllers = controller.NewControllerContainer(cfg, repositories, services)
	middlewares = middleware.NewMiddlewareContainer(services)
}

func main() {
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

	addr := fmt.Sprintf("%s:%s", cfg.ServerConfig.Host, cfg.ServerConfig.Port)

	http.ListenAndServe(addr, nil)
}

func newDB() *sql.DB {
	db, err := sql.Open(
		cfg.DBConfig.DriverName,
		fmt.Sprintf(
			"host=%s port=%d sslmode=%s user=%s password=%s dbname=%s",
			cfg.DBConfig.Host,
			cfg.DBConfig.Port,
			cfg.DBConfig.SSLMode,
			cfg.DBConfig.User,
			cfg.DBConfig.Password,
			cfg.DBConfig.DBName,
		),
	)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db
}
