package model

import (
	"strings"
	"unicode/utf8"

	"github.com/guregu/null"
)

var (
	errCountryIDEmpty          = NewValidationError("id", "идентификатор страны не может быть пустым")
	errCountryNameEmpty        = NewValidationError("name", "название страны не может быть пустым")
	errCountryNameNativeEmpty  = NewValidationError("name_native", "нативное название страны не может быть пустым")
	errCountryEmojiEmpty       = NewValidationError("emoji", "эмодзи страны не может быть пустым")
	errCountryEmojiWrongLength = NewValidationError("emoji", "эмодзи должно содержать один символ")
)

type Country struct {
	ID         int      `json:"id"`
	RegionID   null.Int `json:"region_id"`
	Name       string   `json:"name"`
	NameNative string   `json:"native"`
	Emoji      string   `json:"emoji"`
}

func (c *Country) GetFieldPointers() []interface{} {
	return []interface{}{
		&c.ID,
		&c.RegionID,
		&c.Name,
		&c.NameNative,
		&c.Emoji,
	}
}

func (c *Country) BeforeAdd() {
	c.Name = strings.TrimSpace(c.Name)

	c.NameNative = strings.TrimSpace(c.NameNative)
	if c.NameNative == "" {
		c.NameNative = c.Name
	}

	c.Emoji = strings.TrimSpace(c.Emoji)
}

func (c *Country) Validate() error {
	errors := &ValidationErrors{}

	if c.ID == 0 {
		errors.Append(errCountryIDEmpty)
	}
	if strings.TrimSpace(c.Name) == "" {
		errors.Append(errCountryNameEmpty)
	}
	if strings.TrimSpace(c.NameNative) == "" {
		errors.Append(errCountryNameNativeEmpty)
	}
	if strings.TrimSpace(c.Emoji) == "" {
		errors.Append(errCountryEmojiEmpty)
	}
	if utf8.RuneCountInString(strings.TrimSpace(c.Emoji)) > 2 {
		errors.Append(errCountryEmojiWrongLength)
	}

	if errors.Empty() {
		return nil
	}

	return errors
}
