package router

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router interface {
	Matcher
	HandleFunc(command string, handlerFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update))
}

type Matcher interface {
	Match(identifier string) HandlerFunc
}

type HandlerFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update)

// todo добавить кеширование юзера
type router struct {
	handlers map[string]HandlerFunc
}

func New() Router {
	return &router{handlers: map[string]HandlerFunc{}}
}

func (r *router) HandleFunc(command string, handlerFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update)) {
	r.handlers[command] = HandlerFunc(handlerFunc)
}

func (r *router) Match(identifier string) HandlerFunc {
	for command, handlerFunc := range r.handlers {
		if identifier == command {
			return handlerFunc
		}
	}

	return nil
}
