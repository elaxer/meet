package server

import (
	"errors"
	"meet/internal/app/model"
	"meet/internal/app/repository"
	"net/http"
	"net/url"
	"testing"
)

func TestGetIntQueryParam(t *testing.T) {
	type args struct {
		query     url.Values
		key       string
		byDefault int
		max       int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"A param is under a treshold",
			args{
				query:     url.Values{"param": []string{"15"}},
				key:       "param",
				byDefault: 10,
				max:       20,
			},
			15,
		},
		{
			"A param is above a threshold",
			args{
				query:     url.Values{"param": []string{"21"}},
				key:       "param",
				byDefault: 10,
				max:       20,
			},
			20,
		},
		{
			"A param is an incorrect integer",
			args{
				query:     url.Values{"param": []string{"15abc"}},
				key:       "param",
				byDefault: 10,
				max:       20,
			},
			10,
		},
		{
			"A param is a negative integer",
			args{
				query:     url.Values{"param": []string{"-15"}},
				key:       "param",
				byDefault: 10,
				max:       20,
			},
			10,
		},
		{
			"A param is not transferred",
			args{
				query:     url.Values{},
				key:       "param",
				byDefault: 10,
				max:       20,
			},
			10,
		},
		{
			"Max limit is not specified",
			args{
				query:     url.Values{"param": []string{"999999"}},
				key:       "param",
				byDefault: 0,
				max:       0,
			},
			999999,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetIntQueryParam(tt.args.query, tt.args.key, tt.args.byDefault, tt.args.max); got != tt.want {
				t.Errorf("GetIntQueryParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStatusCode(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"Validation error",
			args{model.NewValidationError("password", "неверный пароль")},
			http.StatusBadRequest,
		},
		{
			"Not Found error",
			args{repository.ErrNotFound},
			http.StatusNotFound,
		},
		{
			"Duplicate error",
			args{repository.ErrDuplicate},
			http.StatusConflict,
		},
		{
			"Unregistered error",
			args{errors.New("unregistered error")},
			http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStatusCode(tt.args.err); got != tt.want {
				t.Errorf("GetStatusCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
