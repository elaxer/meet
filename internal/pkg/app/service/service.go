package service

import (
	"meet/internal/pkg/app"
	"meet/internal/pkg/app/repository"
)

type ServiceContainer struct {
	assessmentService    *AssessmentService
	authService          *AuthService
	fileService          *FileService
	messageService       *MessageService
	photoService         *PhotoService
	questionnaireService *QuestionnaireService
	userService          *UserService
}

func NewServiceContainer(cfg *app.Config, repositories *repository.RepositoryContainer) *ServiceContainer {
	fs := newFileService(cfg.UploadDir)
	as := newAssessmentService(repositories.Assessment(), repositories.Questionnaire())
	return &ServiceContainer{
		assessmentService:    as,
		authService:          newAuthService(cfg.JWTConfig, repositories.User()),
		fileService:          fs,
		messageService:       newMessageService(repositories.Assessment(), repositories.Message(), as),
		photoService:         newPhotoService(repositories.Photo(), repositories.Questionnaire(), fs),
		questionnaireService: newQuestionnaireService(repositories.Questionnaire()),
		userService:          newUserService(repositories.User()),
	}
}

func (sc *ServiceContainer) Assessment() *AssessmentService {
	return sc.assessmentService
}

func (sc *ServiceContainer) Auth() *AuthService {
	return sc.authService
}

func (sc *ServiceContainer) Photo() *PhotoService {
	return sc.photoService
}

func (sc *ServiceContainer) Questionnaire() *QuestionnaireService {
	return sc.questionnaireService
}

func (sc *ServiceContainer) User() *UserService {
	return sc.userService
}
