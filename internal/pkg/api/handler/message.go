package handler

import (
	"encoding/json"
	"meet/internal/pkg/api"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/app/service"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	messageLimitDefault = 100
	messageLimitMax     = 1000
)

type MessageHandler interface {
	GetList(w http.ResponseWriter, r *http.Request)
	UnreadCount(w http.ResponseWriter, r *http.Request)
	Send(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
}

type messageHandler struct {
	messageRepository repository.MessageRepository
	messageService    service.MessageService
}

func NewMessageHandler(messageRepository repository.MessageRepository, messageService service.MessageService) MessageHandler {
	return &messageHandler{messageRepository, messageService}
}

func (mh *messageHandler) GetList(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(api.CtxKeyUser).(*model.User)

	userID, err := api.GetParamInt(mux.Vars(r), "id")
	if err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	query := r.URL.Query()
	limit := api.GetParamQueryInt(query, "limit", messageLimitDefault, messageLimitMax)
	offset := api.GetParamQueryInt(query, "offset", 0, 0)

	ms, err := mh.messageRepository.GetList(model.Direction{FromID: u.ID, ToID: userID}, limit, offset)
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseObject(w, ms, http.StatusOK)
}

func (mh *messageHandler) UnreadCount(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(api.CtxKeyUser).(*model.User)

	count, err := mh.messageRepository.UnreadCount(u.ID)
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))
	}

	response := struct {
		Count int `json:"count"`
	}{Count: count}

	api.ResponseObject(w, response, http.StatusOK)
}

func (mh *messageHandler) Send(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(api.CtxKeyUser).(*model.User)

	m := new(model.Message)
	m.UsersDirection.FromID = u.ID

	if err := json.NewDecoder(r.Body).Decode(m); err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	if err := mh.messageService.Send(m); err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseObject(w, m, http.StatusCreated)
}

func (mh *messageHandler) Read(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(api.CtxKeyUser).(*model.User)

	mID, err := api.GetParamInt(mux.Vars(r), "id")
	if err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	mDTO := new(struct {
		IsReaded bool `json:"is_readed"`
	})
	if err := json.NewDecoder(r.Body).Decode(mDTO); err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}
	if !mDTO.IsReaded {
		api.ResponseError(w, err, http.StatusUnprocessableEntity)

		return
	}

	message, err := mh.messageService.Read(u.ID, mID)
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseObject(w, message, http.StatusOK)
}
