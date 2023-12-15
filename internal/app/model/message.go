package model

type Message struct {
	BaseModel
	Direction Direction `json:"user_direction"`
	Text      string    `json:"text"`
	IsReaded  bool      `json:"is_readed"`
}

// GetFieldPointers реализует интерфейс Model
func (m *Message) GetFieldPointers() []interface{} {
	fields := append(m.BaseModel.GetFieldPointers(), m.Direction.GetFieldPointers()...)

	return append(fields, &m.Text, &m.IsReaded)
}

func (m *Message) Read() {
	m.IsReaded = true
}
