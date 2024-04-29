package router

import (
	"meet/internal/config"
	"meet/internal/pkg/api/handler"
	"meet/internal/pkg/api/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

type configurator struct {
	pathConfig *config.PathConfig

	authorizeMiddleware middleware.AuthorizeMiddleware

	assessmentHandler    handler.AssessmentHandler
	authHandler          handler.AuthHandler
	messageHandler       handler.MessageHandler
	photoHandler         handler.PhotoHandler
	questionnaireHandler handler.QuestionnaireHandler
	swaggerHandler       handler.SwaggerHandler
	userHandler          handler.UserHandler
	dictionaryHandler    handler.DictionaryHandler
}

func NewConfigurator(
	pathConfig *config.PathConfig,

	authorizeMiddleware middleware.AuthorizeMiddleware,

	assessmentHandler handler.AssessmentHandler,
	authHandler handler.AuthHandler,
	messageHandler handler.MessageHandler,
	photoHandler handler.PhotoHandler,
	questionnaireHandler handler.QuestionnaireHandler,
	swaggerHandler handler.SwaggerHandler,
	userHandler handler.UserHandler,
	dictionaryHandler handler.DictionaryHandler,
) *configurator {
	return &configurator{
		pathConfig,
		authorizeMiddleware,
		assessmentHandler,
		authHandler,
		messageHandler,
		photoHandler,
		questionnaireHandler,
		swaggerHandler,
		userHandler,
		dictionaryHandler,
	}
}

func (c *configurator) Configure(r *mux.Router) http.Handler {
	r.StrictSlash(true)

	c.configureFileServer(r)

	apiR := r.PathPrefix("/api/v1").Subrouter()

	c.configureRoutesOpened(apiR)

	apiSecureR := apiR.NewRoute().Subrouter()
	apiSecureR.Use(c.authorizeMiddleware.Authorize)

	c.configureRoutesSecure(apiSecureR)

	return c.configureMiddlewares(r)
}

func (c *configurator) configureMiddlewares(r *mux.Router) http.Handler {
	return middleware.LogMiddleware(
		middleware.CORSMiddleware(
			middleware.ContentLength(
				middleware.FileSize(r),
			),
		),
	)
}

func (c *configurator) configureRoutesOpened(r *mux.Router) {
	r.HandleFunc("/swagger", c.swaggerHandler.Api)

	r.HandleFunc("/users", c.userHandler.Register).Methods(http.MethodPost)
	r.HandleFunc("/authenticate", c.authHandler.Authenticate).Methods(http.MethodPost)

	dictionaryR := r.PathPrefix("/dictionary").Subrouter()
	c.configureRoutesDictionary(dictionaryR)
}

func (c *configurator) configureRoutesDictionary(r *mux.Router) {
	r.HandleFunc("/countries", c.dictionaryHandler.GetCountriesList).Methods(http.MethodGet)
	r.HandleFunc("/countries/{id:[0-9]+}/cities", c.dictionaryHandler.GetCitiesList).Methods(http.MethodGet)
}

func (c *configurator) configureRoutesSecure(r *mux.Router) {
	r.HandleFunc("/users/me", c.userHandler.Me).Methods(http.MethodGet)
	r.HandleFunc("/users/me", c.userHandler.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/passwords/me", c.userHandler.ChangePassword).Methods(http.MethodPut)

	r.HandleFunc("/questionnaires/me", c.questionnaireHandler.Me).Methods(http.MethodGet)
	r.HandleFunc("/questionnaires/me", c.questionnaireHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/questionnaires/me", c.questionnaireHandler.Update).Methods(http.MethodPut)
	r.HandleFunc("/users/me/suggested", c.questionnaireHandler.Suggested).Methods(http.MethodGet)
	r.HandleFunc("/users/me/couples", c.questionnaireHandler.Couples).Methods(http.MethodGet)
	r.HandleFunc("/users/me/assessed", c.questionnaireHandler.Assessed).Methods(http.MethodGet)

	r.HandleFunc("/questionnaires/me/photos/{id:[0-9]+}", c.photoHandler.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/questionnaires/me/photos", c.photoHandler.Upload).Methods(http.MethodPost)

	r.HandleFunc("/assessments", c.assessmentHandler.Assess).Methods(http.MethodPost)

	r.HandleFunc("/users/{id:[0-9]+}/messages", c.messageHandler.List).Methods(http.MethodGet)
	r.HandleFunc("/messages", c.messageHandler.UnreadCount).Methods(http.MethodGet)
	r.HandleFunc("/messages", c.messageHandler.Send).Methods(http.MethodPost)
	r.HandleFunc("/messages/{id:[0-9]+}", c.messageHandler.Read).Methods(http.MethodPut)
}
