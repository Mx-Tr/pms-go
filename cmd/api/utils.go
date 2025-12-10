package main

import (
	"encoding/json"
	"net/http"
	"strings"
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

func (app *Application) UserValidate(w http.ResponseWriter, user *User) []ValidationError {
	var errors []ValidationError

	if !strings.Contains(user.Email, "@") {
		errors = append(errors, ValidationError{"email", "Invalid format"})
	}
	if len(user.Password) < 8 {
		errors = append(errors, ValidationError{"password", "Must be 8+ chars"})
	}

	return errors
}
