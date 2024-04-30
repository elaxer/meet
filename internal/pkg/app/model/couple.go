package model

type Couple struct {
	BaseModel
	UsersDirection Direction `json:"users_direction"`
}

func NewCouple(fromUserID, toUserID int) *Couple {
	couple := new(Couple)
	couple.UsersDirection.FromID = fromUserID
	couple.UsersDirection.ToID = toUserID

	return couple
}

func (c *Couple) GetFieldPointers() []interface{} {
	return append(c.BaseModel.GetFieldPointers(), c.UsersDirection.GetFieldPointers()...)
}

func (c *Couple) BeforeAdd() {
	c.BaseModel.BeforeAdd()
}

func (c *Couple) Validate() error {
	return c.UsersDirection.Validate()
}
