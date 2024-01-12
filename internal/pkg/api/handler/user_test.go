package handler

import (
	"bytes"
	"context"
	"meet/internal/pkg/api"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/app/service"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/guregu/null"
)

var user = model.User{
	BaseModel: model.BaseModel{
		ID:        1,
		CreatedAt: time.Date(2023, 12, 28, 13, 25, 54, 0, time.UTC),
		UpdatedAt: null.TimeFrom(time.Date(2023, 12, 28, 13, 25, 54, 0, time.UTC)),
	},
	Login:        "elaxer",
	PasswordHash: null.StringFrom("!@#$%^&*()_+asdf"),
	IsBlocked:    false,
}

// +integration
func Test_userHandler_Get(t *testing.T) {
	u := user
	ctx := context.WithValue(context.Background(), api.CtxKeyUser, &u)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	NewUserHandler(nil, nil).Get(w, r)

	expected := `{"id":1,"created_at":"2023-12-28T13:25:54Z","login":"elaxer","tg_id":null}`
	got := w.Body.String()
	if got != expected {
		t.Errorf("expected %s got %s", expected, got)
	}
}

// +integration
func Test_userHandler_Register(t *testing.T) {
	requestBody := bytes.NewBuffer([]byte(`{"login": "elaxer", "password": "123456"}`))

	r := httptest.NewRequest(http.MethodPost, "/api/v1/users/me", requestBody)
	w := httptest.NewRecorder()

	ur := repository.NewUserRepository()

	NewUserHandler(nil, service.NewUserService(ur)).Register(w, r)

	expected := http.StatusCreated
	got := w.Result().StatusCode
	if got != expected {
		t.Errorf("expected %d got %d", expected, got)
	}
}

// +integration
func Test_userHandler_ChangePassword(t *testing.T) {
	u := user
	ctx := context.WithValue(context.Background(), api.CtxKeyUser, &u)

	requestBody := bytes.NewBuffer([]byte(`{"login": "elaxer", "password": "654321"}`))

	r := httptest.NewRequest(http.MethodPost, "/api/v1/users/me", requestBody).WithContext(ctx)
	w := httptest.NewRecorder()

	ur := repository.NewUserRepository()
	ur.Add(context.Background(), &u)

	NewUserHandler(nil, service.NewUserService(ur)).ChangePassword(w, r)

	expected := http.StatusNoContent
	got := w.Result().StatusCode
	if got != expected {
		t.Errorf("expected %d got %d", expected, got)
	}
}

// +integration
func Test_userHandler_Delete(t *testing.T) {
	u := user
	ctx := context.WithValue(context.Background(), api.CtxKeyUser, &u)

	requestBody := bytes.NewBuffer([]byte(`{"login": "elaxer", "password": "654321"}`))

	r := httptest.NewRequest(http.MethodDelete, "/api/v1/users/me", requestBody).WithContext(ctx)
	w := httptest.NewRecorder()

	ur := repository.NewUserRepository()
	ur.Add(context.Background(), &u)

	NewUserHandler(nil, service.NewUserService(ur)).ChangePassword(w, r)

	expected := http.StatusNoContent
	got := w.Result().StatusCode
	if got != expected {
		t.Errorf("expected %d got %d", expected, got)
	}
}
