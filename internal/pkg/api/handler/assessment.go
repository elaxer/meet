package handler

import (
	"encoding/json"
	"meet/internal/pkg/api"
	"meet/internal/pkg/app"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/service"
	"net/http"
)

type AssessmentHandler interface {
	Assess(w http.ResponseWriter, r *http.Request)
}

type assessmentHandler struct {
	assessmentService service.AssessmentService
}

func NewAssessmentHandler(assessmentService service.AssessmentService) AssessmentHandler {
	return &assessmentHandler{assessmentService}
}

func (ah *assessmentHandler) Assess(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(app.CtxKeyUser).(*model.User)

	a := new(model.Assessment)
	a.UsersDirection.FromID = u.ID

	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	if err := ah.assessmentService.Assess(a); err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseObject(w, a, http.StatusCreated)
}
