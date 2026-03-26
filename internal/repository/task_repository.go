package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rafaelAmora/internal/model"
)

//Esta camada tem uma responsabilidade: buscar e salvar dados.

type TaskRepository interface {
	FindAll(ctx context.Context) ([]model.Task, error)
	FindByID(ctx context.Context, id string) (model.Task, error)
	Create(ctx context.Context, task model.Task) error
	Update(ctx context.Context, task model.Task) error
	Delete(ctx context.Context, id string) (bool, error)
	Count(ctx context.Context) (int, error)
}

type taskRepository struct {
	db *sqlx.DB // A conexão com o banco, injetada pelo main.go
}

func NewTaskRepository(db *sqlx.DB) TaskRepository {
	return &taskRepository{db: db}
}

// O "receiver" (r *taskRepository) diz: "este método pertence à struct taskRepository".
func (r *taskRepository) FindAll(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task

	err := r.db.SelectContext(ctx, &tasks, "SELECT * FROM tasks ORDER BY created_at DESC")

	return tasks, err
}

func (r *taskRepository) FindByID(ctx context.Context, id string) (model.Task, error) {
	var task model.Task

	err := r.db.GetContext(ctx, &task, "SELECT * FROM tasks WHERE id = $1", id)

	return task, err
}

func (r *taskRepository) Create(ctx context.Context, task model.Task) error {
	query := `
		INSERT INTO tasks (id, title, description, done, created_at)
		VALUES (:id, :title, :description, :done, :created_at)`

	// NamedExecContext usa as tags `db:"..."` para substituir os ":campos"
	// pelos valores reais da struct. Mais legível que usar $1, $2, $3...
	_, err := r.db.NamedExecContext(ctx, query, task)

	return err
}

func (r *taskRepository) Update(ctx context.Context, task model.Task) error {

	query := `UPDATE tasks SET title = :title, description= :description WHERE id = :id`

	_, err := r.db.NamedExecContext(ctx, query, task)

	return err
}

func (r *taskRepository) Delete(ctx context.Context, id string) (bool, error) {

	query := "DELETE FROM tasks WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return false, err
	}

	// RowsAffected retorna quantas linhas foram afetadas.
	// Se for 0, o ID não existia no banco.
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

// Count retorna o total de tarefas cadastradas.
func (r *taskRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM tasks`)
	return count, err
}
