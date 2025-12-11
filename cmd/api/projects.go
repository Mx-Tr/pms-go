package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mx-Tr/pms-go/internal/store"
)

func (app *Application) CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if payload.Name == "" {
		http.Error(w, "Name required", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	project := &store.Project{
		Name:        payload.Name,
		Description: app.GetStringOrEmpty(payload.Description),
		OwnerID:     userID,
	}

	if err := app.store.Projects.Create(r.Context(), project); err != nil {
		fmt.Println("DB Error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.WriteJSON(w, http.StatusOK, project)
}
