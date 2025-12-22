package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Mx-Tr/pms-go/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (app *Application) CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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

	project := &store.Project{
		Name:        payload.Name,
		Description: app.GetStringOrEmpty(payload.Description),
		OwnerId:     userId,
	}

	if err := app.store.Projects.Create(r.Context(), project); err != nil {
		fmt.Println("DB Error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.WriteJSON(w, http.StatusOK, project)
}

func (app *Application) GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO вопрос по поводу названия методов и реалистичного расширения.
	// Не понимаю нужно ли будет для "админа" писать доп правила при получении проектов или
	// нужно будет делать отдельные методы? Узнать.
	ownerId, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := app.store.Projects.GetAll(r.Context(), ownerId)
	if err != nil {
		fmt.Println("DB Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	// Чтобы на фронт не приходил null, пусть туда приходит []
	projects := []store.Project{}

	for rows.Next() {
		var p store.Project

		if err := rows.Scan(&p.Id, &p.Name, &p.Description, &p.OwnerId, &p.CreatedAt); err != nil {
			fmt.Println("Scan Error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		projects = append(projects, p)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("Rows Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.WriteJSON(w, http.StatusOK, projects)
}

func (app *Application) UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	// текущий userId передадим в Update для проверки доступа
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	projectIdStr := chi.URLParam(r, "id")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		http.Error(w, "Invalid projects id", http.StatusBadRequest)
		return
	}

	project, err := app.store.Projects.GetById(r.Context(), projectId, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Project not found", http.StatusNotFound)
		} else {
			fmt.Println("DB Error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	var payload struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		OwnerId     *int    `json:"ownerId"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if app.isAllNil(payload) {
		http.Error(w, "Changes must be", http.StatusBadRequest)
		return
	}

	if payload.Name != nil {
		if *payload.Name == "" {
			http.Error(w, "Name cannot be empty", http.StatusBadRequest)
			return
		}
		project.Name = *payload.Name
	}
	if payload.Description != nil {
		project.Description = *payload.Description
	}
	if payload.OwnerId != nil {
		project.OwnerId = *payload.OwnerId
	}

	err = app.store.Projects.Update(r.Context(), project, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Project not found", http.StatusNotFound)
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" { // код foreign_key_violation
				http.Error(w, "New owner user does not exist", http.StatusBadRequest)
			}
		} else {
			fmt.Println("DB Error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	app.WriteJSON(w, http.StatusOK, project)
}

func (app *Application) DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	projectIdStr := chi.URLParam(r, "id")

	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	err = app.store.Projects.Delete(r.Context(), projectId, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Project not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
