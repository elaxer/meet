package controller

import (
	"encoding/json"
	"meet/internal/api"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/service"
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

	a := new(struct {
		Login    string         `json:"login"`
		Password model.Password `json:"password"`
	})

	err := json.NewDecoder(r.Body).Decode(a)
	if err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	token, err := ac.authService.Authenticate(a.Login, a.Password)
	if err != nil {
		api.ResponseError(w, err, http.StatusUnauthorized)

		return
	}

	api.ResponseObject(w, &responseBody{token}, http.StatusOK)
}
