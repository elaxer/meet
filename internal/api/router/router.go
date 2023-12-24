package router

import (
	"meet/internal/api/controller"
	"meet/internal/api/middleware"

	"github.com/gorilla/mux"
)

func Configure(r *mux.Router, controllers *controller.ControllerContainer, middlewares *middleware.MiddlewareContainer) {
	r.HandleFunc("/swagger/api", controllers.Swagger().Api)

	apiR := r.PathPrefix("/api/v1").Subrouter()

	apiR.HandleFunc("/users", controllers.User().Register).Methods("POST")
	apiR.HandleFunc("/authenticate", controllers.Auth().Authenticate).Methods("POST")

	apiSecureR := apiR.NewRoute().Subrouter()
	apiSecureR.Use(middlewares.Auth().Authorize)

	apiSecureR.HandleFunc("/users/me", controllers.User().Get).Methods("GET")
	apiSecureR.HandleFunc("/users/me", controllers.User().Delete).Methods("DELETE")
	apiSecureR.HandleFunc("/passwords/me", controllers.User().ChangePassword).Methods("PUT")

	apiSecureR.HandleFunc("/questionnaires/me", controllers.Questionnaire().Get).Methods("GET")
	apiSecureR.HandleFunc("/questionnaires/me", controllers.Questionnaire().Create).Methods("POST")
	apiSecureR.HandleFunc("/questionnaires/me", controllers.Questionnaire().Update).Methods("PUT")
	apiSecureR.HandleFunc("/users/me/questionnaires", controllers.Questionnaire().GetList).Methods("GET")
	apiSecureR.HandleFunc("/users/me/couples", controllers.Questionnaire().GetCouples).Methods("GET")

	apiSecureR.HandleFunc("/questionnaires/me/photos/{id:[0-9]+}/content", controllers.Photo().Get).Methods("GET")
	apiSecureR.HandleFunc("/questionnaires/me/photos/{id:[0-9]+}", controllers.Photo().Delete).Methods("DELETE")
	apiSecureR.HandleFunc("/questionnaires/me/photos", controllers.Photo().Upload).Methods("POST")

	apiSecureR.HandleFunc("/assessments", controllers.Assessment().Assess).Methods("POST")

	apiSecureR.HandleFunc("/users/{id:[0-9]+}/messages", controllers.Message().GetList).Methods("GET")
	apiSecureR.HandleFunc("/messages", controllers.Message().Send).Methods("POST")
	apiSecureR.HandleFunc("/messages/{id:[0-9]+}", controllers.Message().Read).Methods("PUT")
}
