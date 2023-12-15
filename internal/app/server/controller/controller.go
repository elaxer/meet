package controller

import (
	"meet/internal/app/repository"
	"meet/internal/app/service"
)

type ControllerContainer struct {
	assessmentController    *assessmentController
	authController          *authController
	messageController       *messageController
	photoController         *photoController
	questionnaireController *questionnaireController
	userController          *userController
}

func NewControllerContainer(repositories *repository.RepositoryContainer, services *service.ServiceContainer) *ControllerContainer {
	return &ControllerContainer{
		assessmentController:    newAssessmentController(services.Assessment()),
		authController:          newAuthController(services.Auth()),
		messageController:       newMessageController(),
		photoController:         newPhotoController(repositories.Photo(), services.Photo()),
		questionnaireController: newQuestionnaireController(repositories.Questionnaire(), services.Questionnaire()),
		userController:          newUserController(repositories.User(), services.User()),
	}
}

func (cc *ControllerContainer) Assessment() *assessmentController {
	return cc.assessmentController
}

func (cc *ControllerContainer) Auth() *authController {
	return cc.authController
}

func (cc *ControllerContainer) Message() *messageController {
	return cc.messageController
}

func (cc *ControllerContainer) Photo() *photoController {
	return cc.photoController
}

func (cc *ControllerContainer) Questionnaire() *questionnaireController {
	return cc.questionnaireController
}

func (cc *ControllerContainer) User() *userController {
	return cc.userController
}
