package handler

import (
	"encoding/json"
	"meet/internal/pkg/api"
	"meet/internal/pkg/app"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/app/service"
	"net/http"
)

type UserHandler interface {
	Me(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	ChangePassword(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	userRepository repository.UserRepository
	userService    service.UserService
}

func NewUserHandler(userRepository repository.UserRepository, userService service.UserService) UserHandler {
	return &userHandler{userRepository, userService}
}

func (uh *userHandler) Me(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(app.CtxKeyUser).(*model.User)

	api.ResponseObject(w, u, http.StatusOK)
}

func (uh *userHandler) Register(w http.ResponseWriter, r *http.Request) {
	uDTO := new(struct {
		Login    string         `json:"login"`
		Password model.Password `json:"password"`
	})

	if err := json.NewDecoder(r.Body).Decode(uDTO); err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	user, err := uh.userService.Register(uDTO.Login, uDTO.Password)
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseObject(w, user, http.StatusCreated)
}

func (uh *userHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(app.CtxKeyUser).(*model.User)

	pDTO := &struct {
		Password model.Password `json:"password"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(pDTO); err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	err := uh.userService.ChangePassword(u, pDTO.Password)
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseEmpty(w, http.StatusNoContent)
}

func (uh *userHandler) Delete(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(app.CtxKeyUser).(*model.User)

	if err := uh.userRepository.Remove(u); err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseObject(w, u, http.StatusOK)
}
