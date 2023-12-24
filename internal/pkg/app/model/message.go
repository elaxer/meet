package model

import "strings"

const messageLenMax = 2048

var (
	ErrMessageTextEmpty   = NewValidationError("text", "текст сообщения не может быть пустым")
	ErrMessageTextTooLong = NewValidationError("text", "длина текста сообщения не должна превышать %d символов", messageLenMax)
)

type Message struct {
	BaseModel
	UsersDirection Direction `json:"users_direction"`
	Text           string    `json:"text"`
	IsReaded       bool      `json:"is_readed"`
}

// GetFieldPointers реализует интерфейс Model
func (m *Message) GetFieldPointers() []interface{} {
	fields := append(m.BaseModel.GetFieldPointers(), m.UsersDirection.GetFieldPointers()...)

	return append(fields, &m.Text, &m.IsReaded)
}

func (m *Message) BeforeAdd() {
	m.BaseModel.BeforeAdd()

	m.Text = strings.TrimSpace(m.Text)
}

func (m *Message) BeforeUpdate() {
	m.BaseModel.BeforeUpdate()

	m.Text = strings.TrimSpace(m.Text)
}

func (m *Message) Validate() error {
	errs := &ValidationErrors{}

	if err := m.UsersDirection.Validate(); err != nil {
		errs.Append(err)
	}
	if strings.TrimSpace(m.Text) == "" {
		errs.Append(ErrMessageTextEmpty)
	}
	if len(m.Text) > messageLenMax {
		errs.Append(ErrMessageTextTooLong)
	}

	if errs.Empty() {
		return nil
	}

	return errs
}

func (m *Message) Read() {
	m.IsReaded = true
}
