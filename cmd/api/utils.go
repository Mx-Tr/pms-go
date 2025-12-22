package main

import (
	"encoding/json"
	"net/http"
	"reflect"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (app *Application) WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func (app *Application) GetStringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (app *Application) isAllNil(v interface{}) bool {
	value := reflect.ValueOf(v)
	for i := 0; i < value.NumField(); i++ {
		if !value.Field(i).IsNil() {
			return false
		}
	}
	return true
}
