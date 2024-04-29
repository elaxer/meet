package handler

import (
	"meet/internal/pkg/api"
	"meet/internal/pkg/app"
	"meet/internal/pkg/app/helper"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/app/service"
	"net/http"

	"github.com/gorilla/mux"
)

type PhotoHandler interface {
	Upload(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type photoHandler struct {
	urlHelper       helper.URLHelper
	photoRepository repository.PhotoRepository
	photoService    service.PhotoService
}

func NewPhotoHandler(urlHelper helper.URLHelper, photoRepository repository.PhotoRepository, photoService service.PhotoService) PhotoHandler {
	return &photoHandler{urlHelper, photoRepository, photoService}
}

func (ph *photoHandler) Upload(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(app.CtxKeyUser).(*model.User)

	file, _, err := r.FormFile("photo")
	if err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}
	defer file.Close()

	photo, err := ph.photoService.Upload(u.ID, file)
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	photo.URL = ph.urlHelper.UploadURL(photo.Path, app.UploadTypeImage, u.ID)

	api.ResponseObject(w, photo, http.StatusCreated)
}

func (ph *photoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(app.CtxKeyUser).(*model.User)

	pID, err := api.GetParamInt(mux.Vars(r), "id")
	if err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	photo, err := ph.photoService.Delete(u.ID, pID)
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	photo.URL = ph.urlHelper.UploadURL(photo.Path, app.UploadTypeImage, u.ID)

	api.ResponseObject(w, photo, http.StatusOK)
}
