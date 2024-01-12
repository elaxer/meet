package handler

import (
	"meet/internal/pkg/api"
	"meet/internal/pkg/app/repository"
	"net/http"

	"github.com/gorilla/mux"
)

type DictionaryHandler interface {
	GetCountriesList(w http.ResponseWriter, r *http.Request)
	GetCitiesList(w http.ResponseWriter, r *http.Request)
}

type dictionaryHandler struct {
	countryRepository repository.CountryRepository
	cityRepository    repository.CityRepository
}

func NewDictionaryHandler(countryRepository repository.CountryRepository, cityRepository repository.CityRepository) DictionaryHandler {
	return &dictionaryHandler{countryRepository, cityRepository}
}

func (dh *dictionaryHandler) GetCountriesList(w http.ResponseWriter, r *http.Request) {
	countries, err := dh.countryRepository.GetAll()
	if err != nil {
		api.ResponseError(w, err, http.StatusInternalServerError)

		return
	}

	api.ResponseObject(w, countries, http.StatusOK)
}

func (dh *dictionaryHandler) GetCitiesList(w http.ResponseWriter, r *http.Request) {
	countryID, err := api.GetParamInt(mux.Vars(r), "id")
	if err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	query := r.URL.Query()
	limit := api.GetParamQueryInt(query, "limit", 100, 0)
	offset := api.GetParamQueryInt(query, "offset", 100, 0)

	cities, err := dh.cityRepository.GetByCountryID(countryID, limit, offset)
	if err != nil {
		api.ResponseError(w, err, http.StatusInternalServerError)

		return
	}

	api.ResponseObject(w, cities, http.StatusOK)
}
