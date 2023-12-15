package model

import "errors"

type Direction struct {
	FromID int `json:"from_id"`
	ToID   int `json:"to_id"`
}

// GetFieldPointers реализует интерфейс Model
func (d *Direction) GetFieldPointers() []interface{} {
	return []interface{}{
		&d.FromID,
		&d.ToID,
	}
}

func (d *Direction) Validate() error {
	if d.FromID == d.ToID {
		return errors.New("идентификаторы отправителя и получателя не должны быть одинаковыми")
	}

	return nil
}

func (d *Direction) NewReversed() Direction {
	return Direction{d.ToID, d.FromID}
}
