package controller

import (
	"encoding/json"
	"meet/internal/app"
	"meet/internal/app/model"
	"meet/internal/app/repository"
	"meet/internal/app/server"
	"meet/internal/app/service"
	"net/http"
)

var (
	questionnaireLimitDefault = app.ListLimitDefault
	questionnaireLimitMax     = app.ListLimitMax
	couplesLimitDefault       = app.ListLimitDefault
	coupleLimitMax            = app.ListLimitMax
)

type questionnaireController struct {
	questionnaireRepository repository.QuestionnaireRepository
	questionnaireService    *service.QuestionnaireService
}

func newQuestionnaireController(
	questionnaireRepository repository.QuestionnaireRepository,
	questionnaireService *service.QuestionnaireService,
) *questionnaireController {
	return &questionnaireController{
		questionnaireRepository: questionnaireRepository,
		questionnaireService:    questionnaireService,
	}
}

func (qc *questionnaireController) Get(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	q, err := qc.questionnaireRepository.GetByUserID(u.ID)
	if err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	server.Response(w, q, http.StatusOK)
}

func (qc *questionnaireController) GetCouples(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	q, err := qc.questionnaireRepository.GetByUserID(u.ID)
	if err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	query := r.URL.Query()
	limit := server.GetIntQueryParam(query, "limit", couplesLimitDefault, coupleLimitMax)
	offset := server.GetIntQueryParam(query, "offset", 0, 0)

	qs, err := qc.questionnaireRepository.GetCouples(q.ID, limit, offset)
	if err != nil {
		server.ResponseError(w, err, http.StatusInternalServerError)

		return
	}

	server.Response(w, qs, http.StatusOK)
}

func (qc *questionnaireController) GetList(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	query := r.URL.Query()
	limit := server.GetIntQueryParam(query, "limit", questionnaireLimitDefault, questionnaireLimitMax)
	offset := server.GetIntQueryParam(query, "offset", 0, 0)

	qs, err := qc.questionnaireService.PickUp(u.ID, limit, offset)
	if err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	server.Response(w, qs, http.StatusOK)
}

func (qc *questionnaireController) Create(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	q := model.NewQuestionnaire()
	q.UserID = u.ID

	if err := json.NewDecoder(r.Body).Decode(q); err != nil {
		server.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	if err := qc.questionnaireService.Add(q); err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	server.Response(w, q, http.StatusCreated)
}

func (qc *questionnaireController) Update(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	q := model.NewQuestionnaire()
	q.UserID = u.ID

	if err := json.NewDecoder(r.Body).Decode(q); err != nil {
		server.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	if err := qc.questionnaireService.Update(q); err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	server.Response(w, q, http.StatusOK)
}
