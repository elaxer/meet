package service

import (
	"errors"
	"fmt"
	"meet/internal/app"
	"meet/internal/app/model"
	"meet/internal/app/repository"
	"time"

	"github.com/golang-jwt/jwt"
)

var ErrFailedLogIn = errors.New("неверный логин или пароль")

type AuthService struct {
	config         *app.Config
	userRepository repository.UserRepository
}

func newAuthService(cfg *app.Config, userRepository repository.UserRepository) *AuthService {
	return &AuthService{
		config:         cfg,
		userRepository: userRepository,
	}
}

func (as *AuthService) Register(login string, password model.Password) error {
	hasUser, err := as.userRepository.HasByLogin(login)
	if err != nil {
		return err
	}
	if hasUser {
		return fmt.Errorf("пользователь с логином \"%s\" уже существует", login)
	}

	if err := password.Validate(); err != nil {
		return err
	}

	u := model.NewUser()
	u.Login = login
	u.PasswordHash, err = password.GetHash()
	if err != nil {
		return err
	}

	err = as.userRepository.Add(u)

	return err
}

func (as *AuthService) Authenticate(login string, password model.Password) (string, error) {
	u, err := as.userRepository.GetByLogin(login)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", ErrFailedLogIn
	}

	if !u.ComparePassword(password) {
		return "", ErrFailedLogIn
	}
	payload := jwt.MapClaims{
		"sub": login,
		"exp": time.Now().Add(time.Second * time.Duration(as.config.JWTConfig.Expire)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString([]byte(as.config.JWTConfig.SecretKey))

	return t, err
}

func (as *AuthService) Authorize(tokenString string) (*model.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неизвестный метод криптографии токена: %v", token.Header["alg"])
		}

		return []byte(as.config.JWTConfig.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("ошибка дешифрования токена")
	}
	if err := claims.Valid(); err != nil {
		return nil, err
	}

	login, ok := claims["sub"]
	if !ok {
		return nil, errors.New("некорректный токен")
	}
	loginString := login.(string)

	u, err := as.userRepository.GetByLogin(loginString)

	return u, err
}
