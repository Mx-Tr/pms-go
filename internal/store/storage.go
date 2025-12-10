package store

import "github.com/jackc/pgx/v5/pgxpool"

type Storage struct {
	Users *UsersRepository
	// Позже тут будут Projects, Tasks...
}

// NewStorage создает экземпляр хранилища
func New(db *pgxpool.Pool) *Storage {
	return &Storage{
		Users: &UsersRepository{db: db},
	}
}
