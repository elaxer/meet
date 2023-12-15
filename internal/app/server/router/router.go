package router

import (
	"meet/internal/app/server/controller"
	"meet/internal/app/server/middleware"

	"github.com/gorilla/mux"
)

func New() *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()

	return s
}

func Configure(r *mux.Router, controllers *controller.ControllerContainer, middlewares *middleware.MiddlewareContainer) {
	r.Use(middleware.ContentLength, middleware.FileSize)

	r.HandleFunc("/users", controllers.User().Register).Methods("POST")
	r.HandleFunc("/authenticate", controllers.Auth().Authenticate).Methods("POST")

	ar := r.NewRoute().Subrouter()
	ar.Use(middlewares.Auth().Authorize)

	ar.HandleFunc("/users/me", controllers.User().Get).Methods("GET")
	ar.HandleFunc("/users/me", controllers.User().Delete).Methods("DELETE")
	ar.HandleFunc("/passwords/me", controllers.User().ChangePassword).Methods("PUT")

	ar.HandleFunc("/questionnaires/me", controllers.Questionnaire().Get).Methods("GET")
	ar.HandleFunc("/questionnaires/me", controllers.Questionnaire().Create).Methods("POST")

	ar.HandleFunc("/questionnaires/me/photos/{id:[0-9]+}", controllers.Photo().Get).Methods("GET")
	ar.HandleFunc("/questionnaires/me/photos/{id:[0-9]+}", controllers.Photo().Delete).Methods("DELETE")
	ar.HandleFunc("/questionnaires/me/photos", controllers.Photo().Upload).Methods("POST")

	ar.HandleFunc("/questionnaires/me", controllers.Questionnaire().Update).Methods("PUT")
	ar.HandleFunc("/users/me/couples", controllers.Questionnaire().GetCouples).Methods("GET")
	ar.HandleFunc("/users/me/questionnaires", controllers.Questionnaire().GetList).Methods("GET")

	ar.HandleFunc("/assessments", controllers.Assessment().Assess).Methods("POST")

	ar.HandleFunc("/users/{id:[0-9]+}/messages", controllers.Message().GetList).Methods("GET")
	ar.HandleFunc("/messages", controllers.Message().Send).Methods("POST")
	ar.HandleFunc("/messages/{id:[0-9]+}", controllers.Message().Read).Methods("PUT")
}
