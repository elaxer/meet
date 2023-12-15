package controller

import (
	"encoding/json"
	"meet/internal/app/model"
	"meet/internal/app/repository"
	"meet/internal/app/server"
	"meet/internal/app/service"
	"net/http"
)

type userController struct {
	userRepository repository.UserRepository
	userService    *service.UserService
}

func newUserController(userRepository repository.UserRepository, userService *service.UserService) *userController {
	return &userController{
		userRepository: userRepository,
		userService:    userService,
	}
}

func (uc *userController) Get(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	server.Response(w, u, http.StatusOK)
}

func (uc *userController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	password := &struct {
		Password model.Password `json:"password"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(password); err != nil {
		server.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	err := uc.userService.ChangePassword(u, password.Password)
	if err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	server.Response(w, nil, http.StatusNoContent)
}

func (uc *userController) Delete(w http.ResponseWriter, r *http.Request) {

}
