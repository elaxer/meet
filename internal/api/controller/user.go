package controller

import (
	"encoding/json"
	"meet/internal/api"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/app/service"
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
	u := r.Context().Value(api.CtxKeyUser).(*model.User)

	api.ResponseObject(w, u, http.StatusOK)
}

func (uc *userController) Register(w http.ResponseWriter, r *http.Request) {
	u := new(struct {
		Login    string         `json:"login"`
		Password model.Password `json:"password"`
	})

	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	user, err := uc.userService.Register(u.Login, u.Password)
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseObject(w, user, http.StatusCreated)
}

func (uc *userController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(api.CtxKeyUser).(*model.User)

	password := &struct {
		Password model.Password `json:"password"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(password); err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	err := uc.userService.ChangePassword(u, password.Password)
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseEmpty(w, http.StatusNoContent)
}

func (uc *userController) Delete(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(api.CtxKeyUser).(*model.User)

	if err := uc.userRepository.Remove(u); err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseObject(w, u, http.StatusOK)
}
