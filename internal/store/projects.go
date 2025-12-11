package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
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

	err := r.db.QueryRow(
		ctx, query, project.Name, project.Description, project.OwnerID,
	).Scan(&project.ID, &project.CreatedAt)

	return err
}

func (r *ProjectsRepository) GetById(ctx context.Context, id int, ownerID int) (*Project, error) {
	project := &Project{}

	query := `SELECT id, name, description, owner_id, created_at FROM projects WHERE id = $1 AND owner_id = $2`

	err := r.db.QueryRow(ctx, query, id, ownerID).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.OwnerID,
		&project.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (r *ProjectsRepository) GetAll(ctx context.Context, ownerID int) (pgx.Rows, error) {
	query := `SELECT id, name, description, owner_id, created_at FROM projects WHERE owner_id = $1`

	rows, err := r.db.Query(ctx, query, ownerID)

	return rows, err
}

func (r *ProjectsRepository) Update(ctx context.Context, project *Project, oldOwnerID int) error {
	query := `UPDATE projects SET name = $1, description = $2, owner_id = $3 WHERE id = $4 and owner_id = $5`
	// TODO нужна ли отдельная ошибка, если нет доступа к проекту? Сейчас возвращается просто Project not found
	// TODO Также ошибка, когда владелец проекта пытается поменять проект на несуществующий ID
	// в консоли возвращается "DB Error: ERROR: insert or update on table "projects" violates foreign key constraint "projects_owner_id_fkey" (SQLSTATE 23503)"

	tag, err := r.db.Exec(
		ctx, query,
		project.Name,
		project.Description,
		project.OwnerID,
		project.ID,
		oldOwnerID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
