package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Mx-Tr/pms-go/internal/store"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
)

var validate = validator.New()

func (app *Application) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload struct {
		Name        string             `json:"name" validate:"required,min=1,max=200"`
		Description *string            `json:"description"`
		Priority    store.TaskPriority `json:"priority" validate:"required,oneof=low medium high urgent"`
		ProjectId   int                `json:"projectId" validate:"required,min=1"`
	}
	// TODO написать ютилити readJSON функцию
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(payload); err != nil {
		app.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if _, err := app.store.Projects.GetById(r.Context(), payload.ProjectId, userId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Project not found", http.StatusBadRequest)
		} else {
			fmt.Println("DB Error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	task := &store.Task{
		Name:        payload.Name,
		Description: app.GetStringOrEmpty(payload.Description),
		Priority:    payload.Priority,
		ProjectId:   payload.ProjectId,
	}

	if err := app.store.Tasks.Create(r.Context(), task); err != nil {
		fmt.Println("DB Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := app.WriteJSON(w, http.StatusOK, task); err != nil {
		log.Printf("Failed to write JSON response: %v", err)
	}
}
