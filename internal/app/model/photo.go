package model

type Photo struct {
	BaseModel
	QuestionnaireID int    `json:"-"`
	Path            string `json:"-"`
}

func (p *Photo) GetFieldPointers() []interface{} {
	return append(p.BaseModel.GetFieldPointers(), &p.QuestionnaireID, &p.Path)
}

func (p *Photo) Validate() error {
	if p.QuestionnaireID == 0 {
		return NewValidationError("questionnaireID", "значение не может быть пустым")
	}
	if p.Path == "" {
		return NewValidationError("path", "значение не может быть пустым")
	}

	return nil
}
