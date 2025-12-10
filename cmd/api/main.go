package main

import (
	"context"
	"fmt"
	"log"
	"project-management-system/internal/config"
	"project-management-system/internal/store"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Загружаем конфиг
	cfg := config.Load()

	fmt.Printf("App starting on port %s\n", cfg.Port)
	fmt.Printf("Connecting to DB: %s\n", cfg.DatabaseURL)

	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("Database if not reachable:", err)
	}

	log.Println("Database connection established")

	// Инициализируем хранилище
	store := store.New(db)
	log.Printf("Storage initialized. Users repo ready: %v", store.Users != nil)

	// Инициализируем приложение
	app := &Application{
		cfg:   cfg,
		store: store,
	}

	// Создаем роутер
	mux := app.Mount()

	// Запуск
	log.Fatal(app.Run(mux))
}
