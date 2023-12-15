package controller

import (
	"encoding/json"
	"meet/internal/app"
	"meet/internal/app/model"
	"meet/internal/app/repository"
	"meet/internal/app/server"
	"meet/internal/app/service"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	messageLimitDefault = app.ListLimitDefault
	messageLimitMax     = app.ListLimitMax
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

	userID, err := server.GetIntParam(mux.Vars(r), "id")
	if err != nil {
		server.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	query := r.URL.Query()
	limit := server.GetIntQueryParam(query, "limit", messageLimitDefault, messageLimitMax)
	offset := server.GetIntQueryParam(query, "offset", 0, 0)

	ms, err := mc.messageRepository.GetList(model.Direction{FromID: u.ID, ToID: userID}, limit, offset)
	if err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	server.Response(w, ms, http.StatusOK)
}

func (mc *messageController) Send(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	m := new(struct {
		ToUserID int    `json:"to_user_id"`
		Text     string `json:"text"`
	})

	if err := json.NewDecoder(r.Body).Decode(m); err != nil {
		server.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	message, err := mc.messageService.Text(model.Direction{FromID: u.ID, ToID: m.ToUserID}, m.Text)
	if err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	server.Response(w, message, http.StatusCreated)
}

func (mc *messageController) Read(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	mID, err := server.GetIntParam(mux.Vars(r), "id")
	if err != nil {
		server.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	m, err := mc.messageService.Read(u.ID, mID)
	if err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	server.Response(w, m, http.StatusOK)
}
