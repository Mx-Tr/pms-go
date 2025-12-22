package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskPriority string
type TaskStatus string

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
	PriorityUrgent TaskPriority = "urgent"
)

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusReview     TaskStatus = "review"
	StatusDone       TaskStatus = "done"
	StatusFrozen     TaskStatus = "frozen"
)

type Task struct {
	Id          int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Priority    TaskPriority `json:"priority"`
	Status      TaskStatus   `json:"status"`
	ProjectId   int          `json:"projectId"`
	CreatedAt   time.Time    `json:"createdAt"`
}

type TasksRepository struct {
	db *pgxpool.Pool
}

func (r *TasksRepository) Create(ctx context.Context, task *Task) error {
	query := `
			INSERT INTO tasks (name, description, priority, project_id) 
			values ($1, $2, $3, $4) 
			RETURNING id, status, created_at
			`

	err := r.db.QueryRow(
		ctx, query,
		task.Name,
		task.Description,
		task.Priority,
		task.ProjectId,
	).Scan(
		&task.Id,
		&task.Status,
		&task.CreatedAt,
	)

	return err
}
