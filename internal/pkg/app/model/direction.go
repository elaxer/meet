package model

var (
	errDirectionIdentifiersEqual = NewValidationError("direction", "идентификаторы отправителя и получателя не должны быть одинаковыми")
)

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
	errs := &ValidationErrors{}
	if d.FromID == d.ToID {
		errs.Append(errDirectionIdentifiersEqual)
	}

	if errs.Empty() {
		return nil
	}

	return errs
}

func (d *Direction) NewReversed() Direction {
	return Direction{d.ToID, d.FromID}
}
