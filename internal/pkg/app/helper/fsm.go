package helper

import (
	"fmt"
)

func NewFSMOccurError(state string) error {
	return fmt.Errorf("невозможно перейти к состоянию %s", state)
}
