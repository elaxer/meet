package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"meet/internal/pkg/app/model"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/h2non/filetype"
)

var (
	errParamNotSpecified = errors.New("параметр не был передан")
)

func ResponseObject(w http.ResponseWriter, v any, statusCode int) {
	if v == nil {
		ResponseEmpty(w, statusCode)

		return
	}

	j, err := json.Marshal(v)
	if err != nil {
		ResponseError(w, err, http.StatusInternalServerError)

		return
	}

	ResponseRaw(w, j, statusCode)
}

func ResponseError(w http.ResponseWriter, err error, statusCode int) {
	log.Println(err)

	switch err.(type) {
	case *model.ValidationError:
		errs := &model.ValidationErrors{}
		errs.Append(err)

		ResponseObject(w, errs, http.StatusUnprocessableEntity)
	case *model.ValidationErrors:
		ResponseObject(w, err, http.StatusUnprocessableEntity)

		return
	}

	ResponseEmpty(w, statusCode)
}

func ResponseRaw(w http.ResponseWriter, bytes []byte, statusCode int) {
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if _, err := w.Write(bytes); err != nil {
		ResponseError(w, err, http.StatusInternalServerError)

		return
	}
}

func ResponseEmpty(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-type", "application/json; charset=utf-8")

	w.WriteHeader(statusCode)
}

func ResponseFile(w http.ResponseWriter, r *http.Request, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		ResponseError(w, err, http.StatusInternalServerError)

		return
	}
	defer file.Close()

	t, err := filetype.MatchReader(file)
	if err != nil {
		ResponseError(w, err, http.StatusInternalServerError)

		return
	}

	fileInfo, err := file.Stat()
	if err != nil {
		ResponseError(w, err, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-type", t.MIME.Value)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	for sc, errs := range errorsMap {
		for _, e := range errs {
			if errors.Is(err, e) {
				return sc
			}
		}
	}

	if _, ok := err.(*model.ValidationError); ok {
		return http.StatusUnprocessableEntity
	}
	if _, ok := err.(*model.ValidationErrors); ok {
		return http.StatusUnprocessableEntity
	}

	return http.StatusInternalServerError
}

func GetParamInt(vars map[string]string, key string) (int, error) {
	v, ok := vars[key]
	if !ok {
		return 0, errParamNotSpecified
	}

	vInt, err := strconv.Atoi(v)

	return vInt, err
}

func GetParamQueryInt(query url.Values, key string, byDefault, max int) int {
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
