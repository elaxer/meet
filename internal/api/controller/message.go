package controller

import (
	"encoding/json"
	"meet/internal/api"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/app/service"
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
	u := r.Context().Value(api.CtxKeyUser).(*model.User)

	userID, err := api.GetParamInt(mux.Vars(r), "id")
	if err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	query := r.URL.Query()
	limit := api.GetParamQueryInt(query, "limit", messageLimitDefault, messageLimitMax)
	offset := api.GetParamQueryInt(query, "offset", 0, 0)

	ms, err := mc.messageRepository.GetList(model.Direction{FromID: u.ID, ToID: userID}, limit, offset)
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseObject(w, ms, http.StatusOK)
}

func (mc *messageController) Send(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(api.CtxKeyUser).(*model.User)

	m := new(model.Message)
	m.UsersDirection.FromID = u.ID

	if err := json.NewDecoder(r.Body).Decode(m); err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	if err := mc.messageService.Text(m); err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseObject(w, m, http.StatusCreated)
}

func (mc *messageController) Read(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(api.CtxKeyUser).(*model.User)

	mID, err := api.GetParamInt(mux.Vars(r), "id")
	if err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	m := new(struct {
		IsReaded bool `json:"is_readed"`
	})
	if err := json.NewDecoder(r.Body).Decode(m); err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}
	if !m.IsReaded {
		api.ResponseError(w, err, http.StatusUnprocessableEntity)

		return
	}

	message, err := mc.messageService.Read(u.ID, mID)
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseObject(w, message, http.StatusOK)
}
