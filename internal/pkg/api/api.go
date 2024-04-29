package api

import (
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/app/service"
	"net/http"
)

var errorsMap = map[int][]error{
	http.StatusBadRequest: {
		errParamNotSpecified,
		service.ErrUsersNotCoupled,
		service.ErrTokenDecoding,
		service.ErrIncorrectCryptoMethod,
	},
	http.StatusUnauthorized:         {service.ErrFailedLogIn},
	http.StatusForbidden:            {service.ErrPhotoUploadLimit, service.ErrQuestionnaireState},
	http.StatusNotFound:             {repository.ErrNotFound},
	http.StatusConflict:             {repository.ErrDuplicate, service.ErrAlreadyAssessed, service.ErrQuestionnairesIncompatible},
	http.StatusUnsupportedMediaType: {service.ErrFileTypeWrong},
}
