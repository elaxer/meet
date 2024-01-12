package router

type compositeRouter struct {
	commandRouter, stateRouter Router
}

func NewCompositeRouter(commandRouter, stateRouter Router) *compositeRouter {
	return &compositeRouter{commandRouter, stateRouter}
}

func (cr *compositeRouter) Match(identifier string) HandlerFunc {
	handler := cr.commandRouter.Match(identifier)
	if handler != nil {
		return handler
	}

	return cr.stateRouter.Match(identifier)
}
