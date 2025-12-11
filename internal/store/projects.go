package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     int       `json:"ownerID"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ProjectsRepository struct {
	db *pgxpool.Pool
}

func (r *ProjectsRepository) Create(ctx context.Context, project *Project) error {
	query := `INSERT INTO projects (name, description, owner_id) values ($1, $2, $3) RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query, project.Name, project.Description, project.OwnerID).Scan(&project.ID, &project.CreatedAt)

	return err
}
