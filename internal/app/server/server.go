package server

import (
	"meet/internal/app/repository"
	"meet/internal/app/service"
	"net/http"
)

type ctxKey int

const (
	CtxKeyUser ctxKey = iota
)

var ErrorsMap = map[int][]error{
	http.StatusBadRequest: {
		ErrParamNotSpecified,
		service.ErrUsersNotCoupled,
		service.ErrTokenDecoding,
		service.ErrIncorrectCryptoMethod,
	},
	http.StatusUnauthorized:         {service.ErrFailedLogIn},
	http.StatusForbidden:            {service.ErrPhotoUploadLimit},
	http.StatusNotFound:             {repository.ErrNotFound},
	http.StatusConflict:             {repository.ErrDuplicate, service.ErrAlreadyAssessed, service.ErrQuestionnairesIncompatible},
	http.StatusUnsupportedMediaType: {service.ErrFileTypeWrong},
}
