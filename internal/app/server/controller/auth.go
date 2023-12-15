package controller

import (
	"encoding/json"
	"meet/internal/app/model"
	"meet/internal/app/server"
	"meet/internal/app/service"
	"net/http"
)

type authController struct {
	authService *service.AuthService
}

func newAuthController(authService *service.AuthService) *authController {
	return &authController{
		authService: authService,
	}
}

func (ac *authController) Authenticate(w http.ResponseWriter, r *http.Request) {
	type responseBody struct {
		AccessToken string `json:"access_token"`
	}

	reqB := new(struct {
		Login    string         `json:"login"`
		Password model.Password `json:"password"`
	})

	err := json.NewDecoder(r.Body).Decode(reqB)
	if err != nil {
		server.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	token, err := ac.authService.Authenticate(reqB.Login, reqB.Password)
	if err != nil {
		server.ResponseError(w, err, http.StatusUnauthorized)

		return
	}

	server.Response(w, &responseBody{token}, http.StatusOK)
}
