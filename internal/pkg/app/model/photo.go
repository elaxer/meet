package model

import "strings"

var (
	errPhotoUnspecifiedQuestionnaire = NewValidationError("questionnaireID", "необходимо указать анкету")
	errPhotoEmptyPath                = NewValidationError("path", "путь не может быть пустым")
)

type Photo struct {
	BaseModel
	QuestionnaireID int    `json:"-"`
	Path            string `json:"-"`
	URL             string `json:"url"`
}

func (p *Photo) GetFieldPointers() []interface{} {
	return append(p.BaseModel.GetFieldPointers(), &p.QuestionnaireID, &p.Path)
}

func (p *Photo) BeforeAdd() {
	p.BaseModel.BeforeAdd()

	p.Path = strings.TrimSpace(p.Path)
	p.Path = strings.Trim(p.Path, "/\\")
}

func (p *Photo) Validate() error {
	errs := &ValidationErrors{}
	if p.QuestionnaireID == 0 {
		errs.Append(errPhotoUnspecifiedQuestionnaire)
	}
	if strings.TrimSpace(p.Path) == "" {
		errs.Append(errPhotoEmptyPath)
	}

	if errs.Empty() {
		return nil
	}

	return errs
}
