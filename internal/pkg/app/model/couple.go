package model

type Couple struct {
	BaseModel
	UsersDirection Direction `json:"users_direction"`
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
