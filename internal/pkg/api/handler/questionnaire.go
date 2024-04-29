package handler

import (
	"context"
	"encoding/json"
	"meet/internal/pkg/api"
	"meet/internal/pkg/app"
	"meet/internal/pkg/app/helper"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/app/service"
	"net/http"
)

const (
	questionnaireLimitDefault = 10
	questionnaireLimitMax     = 100
	couplesLimitDefault       = 10
	couplesLimitMax           = 100
)

type QuestionnaireHandler interface {
	Me(w http.ResponseWriter, r *http.Request)
	Couples(w http.ResponseWriter, r *http.Request)
	Suggested(w http.ResponseWriter, r *http.Request)
	Assessed(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type questionnaireHandler struct {
	urlHelper               helper.URLHelper
	questionnaireRepository repository.QuestionnaireRepository
	questionnaireService    service.QuestionnaireService
}

func NewQuestionnaireHandler(
	urlHelper helper.URLHelper,
	questionnaireRepository repository.QuestionnaireRepository,
	questionnaireService service.QuestionnaireService,
) QuestionnaireHandler {
	return &questionnaireHandler{urlHelper, questionnaireRepository, questionnaireService}
}

func (qh *questionnaireHandler) Me(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(app.CtxKeyUser).(*model.User)

	q, err := qh.questionnaireRepository.GetByUserID(u.ID)
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	qh.urlHelper.SetQuestionnairePhotos(q)

	api.ResponseObject(w, q, http.StatusOK)
}

func (qh *questionnaireHandler) Couples(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(app.CtxKeyUser).(*model.User)

	query := r.URL.Query()
	limit := api.GetParamQueryInt(query, "limit", couplesLimitDefault, couplesLimitMax)
	offset := api.GetParamQueryInt(query, "offset", 0, 0)

	qs, err := qh.questionnaireRepository.Couples(u.ID, limit, offset)
	if err != nil {
		api.ResponseError(w, err, http.StatusInternalServerError)

		return
	}

	qh.urlHelper.SetQuestionnairePhotos(qs...)

	api.ResponseObject(w, qs, http.StatusOK)
}

func (qh *questionnaireHandler) Suggested(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(app.CtxKeyUser).(*model.User)

	query := r.URL.Query()
	limit := api.GetParamQueryInt(query, "limit", questionnaireLimitDefault, questionnaireLimitMax)
	offset := api.GetParamQueryInt(query, "offset", 0, 0)

	qs, err := qh.questionnaireService.Suggested(u.ID, limit, offset)
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	qh.urlHelper.SetQuestionnairePhotos(qs...)

	api.ResponseObject(w, qs, http.StatusOK)
}

func (qh *questionnaireHandler) Assessed(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(app.CtxKeyUser).(*model.User)

	query := r.URL.Query()
	limit := api.GetParamQueryInt(query, "limit", questionnaireLimitDefault, questionnaireLimitMax)
	offset := api.GetParamQueryInt(query, "offset", 0, 0)

	qs, err := qh.questionnaireRepository.Assessed(u.ID, limit, offset)
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	qh.urlHelper.SetQuestionnairePhotos(qs...)

	api.ResponseObject(w, qs, http.StatusOK)
}

func (qh *questionnaireHandler) Create(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(app.CtxKeyUser).(*model.User)

	q := model.NewQuestionnaireEmpty()
	q.UserID = u.ID

	if err := json.NewDecoder(r.Body).Decode(q); err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	if err := qh.questionnaireService.Add(context.Background(), q); err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	qh.urlHelper.SetQuestionnairePhotos(q)

	api.ResponseObject(w, q, http.StatusCreated)
}

func (qh *questionnaireHandler) Update(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(app.CtxKeyUser).(*model.User)

	q := new(model.Questionnaire)
	q.UserID = u.ID

	if err := json.NewDecoder(r.Body).Decode(q); err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	if err := qh.questionnaireService.Update(q); err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	qh.urlHelper.SetQuestionnairePhotos(q)

	api.ResponseObject(w, q, http.StatusOK)
}
