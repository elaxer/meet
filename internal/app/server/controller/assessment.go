package controller

import (
	"encoding/json"
	"meet/internal/app/model"
	"meet/internal/app/server"
	"meet/internal/app/service"
	"net/http"

	"github.com/guregu/null"
)

type assessmentController struct {
	assessmentService *service.AssessmentService
}

func newAssessmentController(assessmentService *service.AssessmentService) *assessmentController {
	return &assessmentController{assessmentService}
}

func (ac *assessmentController) Assess(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	a := new(struct {
		QuestionnaireID int            `json:"questionnaire_id"`
		Decision        model.Decision `json:"decision"`
		Message         null.String    `json:"message"`
		IsMutual        bool           `json:"is_mutual"`
	})

	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		server.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	isMutual, err := ac.assessmentService.Assess(u.ID, a.QuestionnaireID, a.Decision, a.Message)
	if err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	a.IsMutual = isMutual

	server.Response(w, a, http.StatusCreated)
}
