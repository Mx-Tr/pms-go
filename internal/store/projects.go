package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Project struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerId     int       `json:"ownerId"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ProjectsRepository struct {
	db *pgxpool.Pool
}

func (r *ProjectsRepository) Create(ctx context.Context, project *Project) error {
	query := `INSERT INTO projects (name, description, owner_id) values ($1, $2, $3) RETURNING id, created_at`

	err := r.db.QueryRow(
		ctx, query, project.Name, project.Description, project.OwnerId,
	).Scan(&project.Id, &project.CreatedAt)

	return err
}

func (r *ProjectsRepository) GetById(ctx context.Context, id int, ownerId int) (*Project, error) {
	project := &Project{}

	query := `SELECT id, name, description, owner_id, created_at FROM projects WHERE id = $1 AND owner_id = $2`

	err := r.db.QueryRow(ctx, query, id, ownerId).Scan(
		&project.Id,
		&project.Name,
		&project.Description,
		&project.OwnerId,
		&project.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (r *ProjectsRepository) GetAll(ctx context.Context, ownerId int) (pgx.Rows, error) {
	query := `SELECT id, name, description, owner_id, created_at FROM projects WHERE owner_id = $1`

	rows, err := r.db.Query(ctx, query, ownerId)

	return rows, err
}

func (r *ProjectsRepository) Update(ctx context.Context, project *Project, oldOwnerId int) error {
	query := `UPDATE projects SET name = $1, description = $2, owner_id = $3 WHERE id = $4 and owner_id = $5`
	// TODO Также ошибка, когда владелец проекта пытается поменять проект на несуществующий Id
	// в консоли возвращается "DB Error: ERROR: insert or update on table "projects" violates foreign key constraint "projects_owner_id_fkey" (SQLSTATE 23503)"

	tag, err := r.db.Exec(
		ctx, query,
		project.Name,
		project.Description,
		project.OwnerId,
		project.Id,
		oldOwnerId,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *ProjectsRepository) Delete(ctx context.Context, projectId int, ownerId int) error {
	query := `DELETE FROM projects WHERE id = $1 and owner_id = $2`

	tag, err := r.db.Exec(ctx, query, projectId, ownerId)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
