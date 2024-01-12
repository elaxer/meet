package handler

import (
	"encoding/json"
	"meet/internal/pkg/api"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/service"
	"net/http"
)

type AuthHandler interface {
	Authenticate(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{authService}
}

func (ah *authHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	type responseBody struct {
		AccessToken string `json:"access_token"`
	}

	a := new(struct {
		Login    string         `json:"login"`
		Password model.Password `json:"password"`
	})

	err := json.NewDecoder(r.Body).Decode(a)
	if err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	token, err := ah.authService.Authenticate(a.Login, a.Password)
	if err != nil {
		api.ResponseError(w, err, http.StatusUnauthorized)

		return
	}

	api.ResponseObject(w, &responseBody{token}, http.StatusOK)
}
