package service

import (
	"context"
	"errors"
	"meet/internal/config"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrFailedLogIn           = errors.New("неверный логин или пароль")
	ErrTokenDecoding         = errors.New("ошибка дешифрования токена")
	ErrIncorrectCryptoMethod = errors.New("некорректный метод криптографии токена")
)

type AuthService interface {
	Authenticate(login string, password model.Password) (string, error)
	Authorize(tokenString string) (*model.User, error)
}

type authService struct {
	jwtConfig           *config.JWTConfig
	userDBRepository    repository.UserRepository
	userRedisRepository repository.UserRepository
	userService         UserService
}

func NewAuthService(
	jwtConfig *config.JWTConfig,
	userDBRepository repository.UserRepository,
	userRedisRepository repository.UserRepository,
	userService UserService,
) AuthService {
	return &authService{jwtConfig, userDBRepository, userRedisRepository, userService}
}

func (as *authService) Authenticate(login string, password model.Password) (string, error) {
	u, err := as.userDBRepository.GetByLogin(login)
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
		"exp": time.Now().Add(time.Second * time.Duration(as.jwtConfig.Expire)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString([]byte(as.jwtConfig.SecretKey))
	if err != nil {
		return "", err
	}

	if err := as.userRedisRepository.Add(context.Background(), u); err != nil {
		return "", err
	}

	return t, nil
}

func (as *authService) Authorize(tokenString string) (*model.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrIncorrectCryptoMethod
		}

		return []byte(as.jwtConfig.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrTokenDecoding
	}
	if err := claims.Valid(); err != nil {
		return nil, err
	}

	login, ok := claims["sub"]
	if !ok {
		return nil, ErrTokenDecoding
	}
	loginString := login.(string)

	return as.userService.GetByLogin(loginString)
}
