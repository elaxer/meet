package main

import (
	"database/sql"
	"fmt"
	"log"
	"meet/internal/app"
	"meet/internal/app/repository"
	"meet/internal/app/server/controller"
	"meet/internal/app/server/middleware"
	"meet/internal/app/server/router"
	"meet/internal/app/service"
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
		log.Fatal("Error loading .env file")
	}

	cfg, err = app.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := newDB()
	if err != nil {
		log.Fatal(err)
	}

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
	http.ListenAndServe(cfg.ServerConfig.Port, nil)
}

func newDB() (*sql.DB, error) {
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
		return nil, err
	}

	err = db.Ping()

	return db, err
}
