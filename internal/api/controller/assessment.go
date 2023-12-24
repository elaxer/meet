package controller

import (
	"encoding/json"
	"meet/internal/api"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/service"
	"net/http"
	"time"
)

type assessmentController struct {
	assessmentService *service.AssessmentService
}

func newAssessmentController(assessmentService *service.AssessmentService) *assessmentController {
	return &assessmentController{assessmentService}
}

func (ac *assessmentController) Assess(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(api.CtxKeyUser).(*model.User)

	a := new(model.Assessment)
	a.UsersDirection.FromID = u.ID

	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	if err := ac.assessmentService.Assess(a, time.Now()); err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseObject(w, a, http.StatusCreated)
}
