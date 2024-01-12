package model

import "strings"

var (
	errCityIDEmpty        = NewValidationError("id", "идентификатор города не может быть пустым")
	errCityCountryIDEmpty = NewValidationError("country_id", "идентификатор страны города не может быть пустым")
	errCityNameEmpty      = NewValidationError("name", "название города не может быть пустым")

	errLatitudeRange  = NewValidationError("latitude", "широта должна быть не менее -90 и не более 90 градусов")
	errLongitudeRange = NewValidationError("longitude", "долгода должна быть не менее -180 и не более 180 градусов")
)

type City struct {
	ID        int    `json:"id"`
	CountryID int    `json:"country_id,omitempty"`
	Name      string `json:"name"`
	Coordinate
}

func (c *City) GetFieldPointers() []interface{} {
	return []interface{}{&c.ID, &c.CountryID, &c.Name, &c.Coordinate.Latitude, &c.Coordinate.Longitude}
}

func (c *City) BeforeAdd() {
	c.Name = strings.TrimSpace(c.Name)
}

func (c *City) Validate() error {
	errors := &ValidationErrors{}

	if c.ID == 0 {
		errors.Append(errCityIDEmpty)
	}
	if c.CountryID == 0 {
		errors.Append(errCityCountryIDEmpty)
	}
	if strings.TrimSpace(c.Name) == "" {
		errors.Append(errCityNameEmpty)
	}

	if err := c.Coordinate.Validate(); err != nil {
		errors.Append(err)
	}

	if errors.Empty() {
		return nil
	}

	return errors
}

type Coordinate struct {
	Latitude  float64 `json:"latitude,string"`
	Longitude float64 `json:"longitude,string"`
}

func (c *Coordinate) Validate() error {
	errors := &ValidationErrors{}

	if c.Latitude < -90 || c.Latitude > 90 {
		errors.Append(errLatitudeRange)
	}

	if c.Longitude < -180 || c.Longitude > 180 {
		errors.Append(errLongitudeRange)
	}

	if errors.Empty() {
		return nil
	}

	return errors
}
