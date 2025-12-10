package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // При ответе клиенту эта строка не вернётся
	CreatedAt    time.Time `json:"created_at"`
}

// UserRepository отвечает за работу с таблицей users
type UsersRepository struct {
	db *pgxpool.Pool
}

// Create добавляет пользователя в базу
func (r *UsersRepository) Create(ctx context.Context, user *User) error {
	query := `
			INSERT INTO users (email, password_hash)
			values ($1, $2)
			RETURNING id, created_at
	`

	// Мы используем данные из структуры user, чтобы заполнить SQL
	// И сразу записываем в user полученные id и время
	err := r.db.QueryRow(ctx, query, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
	return err
}

func (r *UsersRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}

	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
	//  QueryRow вернет ошибку pgx.ErrNoRows, если юзера нет
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	return user, err
}

func (r *UsersRepository) GetById(ctx context.Context, id int) (*User, error) {
	user := &User{}

	query := `SELECT id, email, created_at FROM users WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
	)

	return user, err
}
