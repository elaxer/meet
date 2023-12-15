package controller

import (
	"meet/internal/app/model"
	"meet/internal/app/repository"
	"meet/internal/app/server"
	"meet/internal/app/service"
	"net/http"

	"github.com/gorilla/mux"
)

type photoController struct {
	photoRepository repository.PhotoRepository
	photoService    *service.PhotoService
}

func newPhotoController(photoRepository repository.PhotoRepository, photoService *service.PhotoService) *photoController {
	return &photoController{
		photoRepository: photoRepository,
		photoService:    photoService,
	}
}

func (pc *photoController) Upload(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	file, _, err := r.FormFile("file")
	if err != nil {
		server.ResponseError(w, err, http.StatusBadRequest)

		return
	}
	defer file.Close()

	photo, err := pc.photoService.Upload(u.ID, file)
	if err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	server.Response(w, photo, http.StatusCreated)
}

func (pc *photoController) Get(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	pID, err := server.GetIntParam(mux.Vars(r), "id")
	if err != nil {
		server.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	path, err := pc.photoService.GetPath(u.ID, pID)
	if err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	http.ServeFile(w, r, path)
}

func (pc *photoController) Delete(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(server.CtxKeyUser).(*model.User)

	pID, err := server.GetIntParam(mux.Vars(r), "id")
	if err != nil {
		server.ResponseError(w, err, http.StatusBadRequest)

		return
	}

	p, err := pc.photoService.Delete(u.ID, pID)
	if err != nil {
		server.ResponseError(w, err, server.GetStatusCode(err))

		return
	}

	server.Response(w, p, http.StatusOK)
}
