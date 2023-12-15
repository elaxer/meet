package model

import "strings"

var messageLenMax = 2048

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

func (m *Message) Validate() error {
	if err := m.UsersDirection.Validate(); err != nil {
		return err
	}

	if strings.TrimSpace(m.Text) == "" {
		return NewValidationError("text", "сообщение не может быть пустым")
	}

	if len(m.Text) > messageLenMax {
		return NewValidationError("text", "длина сообщения не должна превышать %d символов", messageLenMax)
	}

	return nil
}

func (m *Message) Read() {
	m.IsReaded = true
}
