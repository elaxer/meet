package controller

import (
	"encoding/json"
	"meet/internal/api"
	"meet/internal/pkg/app"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/app/service"
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
	u := r.Context().Value(api.CtxKeyUser).(*model.User)

	q, err := qc.questionnaireRepository.GetByUserID(u.ID)
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseObject(w, q, http.StatusOK)
}

func (qc *questionnaireController) GetCouples(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(api.CtxKeyUser).(*model.User)

	q, err := qc.questionnaireRepository.GetByUserID(u.ID)
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	query := r.URL.Query()
	limit := api.GetParamQueryInt(query, "limit", couplesLimitDefault, couplesLimitMax)
	offset := api.GetParamQueryInt(query, "offset", 0, 0)

	qs, err := qc.questionnaireRepository.Couples(q.ID, limit, offset)
	if err != nil {
		api.ResponseError(w, err, http.StatusInternalServerError)

		return
	}

	api.ResponseObject(w, qs, http.StatusOK)
}

func (qc *questionnaireController) GetList(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(api.CtxKeyUser).(*model.User)

	query := r.URL.Query()
	limit := api.GetParamQueryInt(query, "limit", questionnaireLimitDefault, questionnaireLimitMax)
	offset := api.GetParamQueryInt(query, "offset", 0, 0)

	qs, err := qc.questionnaireService.PickUp(u.ID, limit, offset, time.Now())
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseObject(w, qs, http.StatusOK)
}

func (qc *questionnaireController) Create(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(api.CtxKeyUser).(*model.User)

	q := new(model.Questionnaire)
	q.UserID = u.ID

	if err := json.NewDecoder(r.Body).Decode(q); err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	if err := qc.questionnaireService.Add(q); err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseObject(w, q, http.StatusCreated)
}

func (qc *questionnaireController) Update(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(api.CtxKeyUser).(*model.User)

	q := new(model.Questionnaire)
	q.UserID = u.ID

	if err := json.NewDecoder(r.Body).Decode(q); err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	if err := qc.questionnaireService.Update(q); err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseObject(w, q, http.StatusOK)
}
