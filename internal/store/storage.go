package store

import "github.com/jackc/pgx/v5/pgxpool"

type Storage struct {
	Users    *UsersRepository
	Projects *ProjectsRepository
	Tasks    *TasksRepository
	// Позже тут будут Projects, Tasks...
}

// NewStorage создает экземпляр хранилища
func New(db *pgxpool.Pool) *Storage {
	return &Storage{
		Users:    &UsersRepository{db: db},
		Projects: &ProjectsRepository{db: db},
		Tasks:    &TasksRepository{db: db},
	}
}
