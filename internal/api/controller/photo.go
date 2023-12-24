package controller

import (
	"meet/internal/api"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/app/service"
	"net/http"

	"github.com/gorilla/mux"
)

type photoController struct {
	photoRepository repository.PhotoRepository
	photoService    *service.PhotoService
}

func newPhotoController(photoRepository repository.PhotoRepository, photoService *service.PhotoService) *photoController {
	return &photoController{photoRepository, photoService}
}

func (pc *photoController) Upload(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(api.CtxKeyUser).(*model.User)

	file, _, err := r.FormFile("photo")
	if err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}
	defer file.Close()

	photo, err := pc.photoService.Upload(u.ID, file)
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseObject(w, photo, http.StatusCreated)
}

func (pc *photoController) Get(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(api.CtxKeyUser).(*model.User)

	pID, err := api.GetParamInt(mux.Vars(r), "id")
	if err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	fp, err := pc.photoService.GetPath(u.ID, pID)
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseFile(w, r, fp)
}

func (pc *photoController) Delete(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(api.CtxKeyUser).(*model.User)

	pID, err := api.GetParamInt(mux.Vars(r), "id")
	if err != nil {
		api.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	p, err := pc.photoService.Delete(u.ID, pID)
	if err != nil {
		api.ResponseError(w, err, api.GetStatusCode(err))

		return
	}

	api.ResponseObject(w, p, http.StatusOK)
}
