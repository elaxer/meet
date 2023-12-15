package server

import (
	"encoding/json"
	"errors"
	"log"
	"meet/internal/app/model"
	"meet/internal/app/repository"
	"meet/internal/app/service"
	"net/http"
	"net/url"
	"strconv"
)

func Response(w http.ResponseWriter, v any, statusCode int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)

	if v == nil {
		return
	}

	j, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
	}

	_, err = w.Write(j)
	if err != nil {
		log.Println(err)
	}
}

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if _, ok := err.(*model.ValidationError); ok {
		return http.StatusBadRequest
	}
	if errors.Is(err, repository.ErrDuplicate) {
		return http.StatusConflict
	}
	if errors.Is(err, repository.ErrNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, service.ErrFileTypeWrong) {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}

func ResponseError(w http.ResponseWriter, err error, statusCode int) {
	log.Println(err)

	Response(w, nil, statusCode)
}

func GetIntParam(vars map[string]string, key string) (int, error) {
	v, ok := vars[key]
	if !ok {
		return 0, errors.New("параметр не был передан")
	}

	vInt, err := strconv.Atoi(v)

	return vInt, err
}

func GetIntQueryParam(query url.Values, key string, byDefault, max int) int {
	v := query.Get(key)
	if v == "" {
		return byDefault
	}

	vInt, err := strconv.Atoi(v)
	if err != nil {
		return byDefault
	}

	if vInt < 0 {
		return byDefault
	}

	if max != 0 && vInt > max {
		return max
	}

	return vInt
}
