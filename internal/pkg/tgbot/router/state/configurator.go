package state

import (
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/tgbot/command"
	"meet/internal/pkg/tgbot/router"
)

type configurator struct {
	qc command.QuestionnaireCommand
}

func NewConfigurator(qc command.QuestionnaireCommand) router.RouterConfigurator {
	return &configurator{qc}
}

func (src *configurator) Configure(r router.Router) {
	r.HandleFunc(model.StateQuestionnaireFillingName, src.qc.FillName)
	r.HandleFunc(model.StateQuestionnaireFillingBirthDate, src.qc.FillBirthDate)
	r.HandleFunc(model.StateQuestionnaireFillingGender, src.qc.FillGender)
}
