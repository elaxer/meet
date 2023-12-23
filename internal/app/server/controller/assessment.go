package controller

import (
	"encoding/json"
	"meet/internal/app/model"
	"meet/internal/app/server"
	"meet/internal/app/service"
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
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	a := new(model.Assessment)
	a.UsersDirection.FromID = u.ID

	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		server.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	if err := ac.assessmentService.Assess(a, time.Now()); err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	server.ResponseObject(w, a, http.StatusCreated)
}
