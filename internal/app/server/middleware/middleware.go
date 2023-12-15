package middleware

import "meet/internal/app/service"

type MiddlewareContainer struct {
	authMiddleware *authMiddleware
}

func NewMiddlewareContainer(services *service.ServiceContainer) *MiddlewareContainer {
	return &MiddlewareContainer{
		authMiddleware: newAuthMiddleware(services.Auth()),
	}
}

func (mc *MiddlewareContainer) Auth() *authMiddleware {
	return mc.authMiddleware
}
