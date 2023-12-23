package controller

import (
	"encoding/json"
	"meet/internal/app/model"
	"meet/internal/app/repository"
	"meet/internal/app/server"
	"meet/internal/app/service"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	messageLimitDefault = 100
	messageLimitMax     = 1000
)

type messageController struct {
	messageRepository repository.MessageRepository
	messageService    *service.MessageService
}

func newMessageController() *messageController {
	return &messageController{}
}

func (mc *messageController) GetList(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	userID, err := server.GetParamInt(mux.Vars(r), "id")
	if err != nil {
		server.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	query := r.URL.Query()
	limit := server.GetParamQueryInt(query, "limit", messageLimitDefault, messageLimitMax)
	offset := server.GetParamQueryInt(query, "offset", 0, 0)

	ms, err := mc.messageRepository.GetList(model.Direction{FromID: u.ID, ToID: userID}, limit, offset)
	if err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	server.ResponseObject(w, ms, http.StatusOK)
}

func (mc *messageController) Send(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	m := new(model.Message)
	m.UsersDirection.FromID = u.ID

	if err := json.NewDecoder(r.Body).Decode(m); err != nil {
		server.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	if err := mc.messageService.Text(m); err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	server.ResponseObject(w, m, http.StatusCreated)
}

func (mc *messageController) Read(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	mID, err := server.GetParamInt(mux.Vars(r), "id")
	if err != nil {
		server.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	m := new(struct {
		IsReaded bool `json:"is_readed"`
	})
	if err := json.NewDecoder(r.Body).Decode(m); err != nil {
		server.ResponseError(w, err, http.StatusBadRequest)

		return
	}
	if !m.IsReaded {
		server.ResponseError(w, err, http.StatusUnprocessableEntity)

		return
	}

	message, err := mc.messageService.Read(u.ID, mID)
	if err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	server.ResponseObject(w, message, http.StatusOK)
}
