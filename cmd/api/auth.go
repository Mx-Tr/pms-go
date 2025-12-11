package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Mx-Tr/pms-go/internal/store"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Декодируем запрос
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Валидация для возвращения массива ошибок
	type ValidationError struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}
	var errors []ValidationError
	if !strings.Contains(payload.Email, "@") {
		errors = append(errors, ValidationError{"email", "Invalid format"})
	}
	if len(payload.Password) < 8 {
		errors = append(errors, ValidationError{"password", "Must be 8+ chars"})
	}

	if len(errors) > 0 {
		app.WriteJSON(w, http.StatusBadRequest, errors)
		return
	}

	// 2. Хэшируем пароль
	// Cost - это сложность шифрования
	// bcrypt.DefaultCost (10) - золотая середина между скоростью и безопасностью
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 3. Создаем модель User
	user := &store.User{
		Email:        payload.Email,
		PasswordHash: string(hashedPassword), // Сохраняем хэш, а не пароль
	}

	// 4. Сохраняем в базу

	if err := app.store.Users.Create(r.Context(), user); err != nil {
		// TODO Добавить проверку на занятый email
		fmt.Println("DB Error:", err)
		http.Error(w, "Could not create user (email must be unique)", http.StatusInternalServerError)
		return
	}

	// 5. Ответ на клиент
	app.WriteJSON(w, http.StatusCreated, user)
}

func (app *Application) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !strings.Contains(payload.Email, "@") || len(payload.Password) < 8 {
		http.Error(w, "Invalid Email or Password", http.StatusBadRequest)
		return
	}

	// 1. Поиск юзера
	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Invalid Email or Password", http.StatusBadRequest)
			return
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	// 2. Сравнение паролей через CompareHashAndPassword
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		http.Error(w, "Invalid Email or Password", http.StatusBadRequest)
		return
	}

	// 3. Генерируем JWT
	// Claims это данные, зашитые в токен (payload)
	claims := jwt.MapClaims{
		"sub": user.ID,                               // Subject (Обычно ID пользователя)
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Expire (срок жизни - 24 часа)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Секретный ключ для подписи из .env
	tokenString, err := token.SignedString([]byte(app.cfg.JWTSecret))
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// 4. Отдаем токен
	app.WriteJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

func (app *Application) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	// Достаем ID из контекста, который положил туда наш Middleware
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := app.store.Users.GetById(r.Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if err := app.WriteJSON(w, http.StatusOK, user); err != nil {
		log.Printf("WriteJSON error: %v", err)
		return
	}

}
