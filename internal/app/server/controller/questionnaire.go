package controller

import (
	"encoding/json"
	"meet/internal/app"
	"meet/internal/app/model"
	"meet/internal/app/repository"
	"meet/internal/app/server"
	"meet/internal/app/service"
	"net/http"
	"time"
)

var (
	questionnaireLimitDefault = app.ListLimitDefault
	questionnaireLimitMax     = app.ListLimitMax
	couplesLimitDefault       = app.ListLimitDefault
	couplesLimitMax           = app.ListLimitMax
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

	server.ResponseObject(w, q, http.StatusOK)
}

func (qc *questionnaireController) GetCouples(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	q, err := qc.questionnaireRepository.GetByUserID(u.ID)
	if err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	query := r.URL.Query()
	limit := server.GetParamQueryInt(query, "limit", couplesLimitDefault, couplesLimitMax)
	offset := server.GetParamQueryInt(query, "offset", 0, 0)

	qs, err := qc.questionnaireRepository.Couples(q.ID, limit, offset)
	if err != nil {
		server.ResponseError(w, err, http.StatusInternalServerError)

		return
	}

	server.ResponseObject(w, qs, http.StatusOK)
}

func (qc *questionnaireController) GetList(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	query := r.URL.Query()
	limit := server.GetParamQueryInt(query, "limit", questionnaireLimitDefault, questionnaireLimitMax)
	offset := server.GetParamQueryInt(query, "offset", 0, 0)

	qs, err := qc.questionnaireService.PickUp(u.ID, limit, offset, time.Now())
	if err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	server.ResponseObject(w, qs, http.StatusOK)
}

func (qc *questionnaireController) Create(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	q := new(model.Questionnaire)
	q.UserID = u.ID

	if err := json.NewDecoder(r.Body).Decode(q); err != nil {
		server.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	if err := qc.questionnaireService.Add(q); err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	server.ResponseObject(w, q, http.StatusCreated)
}

func (qc *questionnaireController) Update(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	q := new(model.Questionnaire)
	q.UserID = u.ID

	if err := json.NewDecoder(r.Body).Decode(q); err != nil {
		server.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	if err := qc.questionnaireService.Update(q); err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	server.ResponseObject(w, q, http.StatusOK)
}
