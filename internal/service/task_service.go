package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rafaelAmora/internal/model"
	"github.com/rafaelAmora/internal/repository"
)

// Esta camada contém as REGRAS DE NEGÓCIO.
// Ela fica entre o Controller (que lida com HTTP) e o Repository (banco de dados).

// Erros de negócio — erros com SIGNIFICADO para o domínio da aplicação.
var (
	ErrTaskNotFound  = errors.New("tarefa não encontrada")
	ErrTitleTooShort = errors.New("o título precisa ter pelo menos 3 caracteres")
	ErrTitleEmpty    = errors.New("o título não pode estar vazio")
	ErrDescEmpty     = errors.New("a descrição não pode estar vazia")
)

type TaskService interface {
	List(ctx context.Context) ([]model.Task, error)
	GetByID(ctx context.Context, id string) (model.Task, error)
	Create(ctx context.Context, input model.CreateTaskInput) (model.Task, error)
	Update(ctx context.Context, id string, input model.UpdateTaskInput) (model.Task, error)
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}

// taskService é a implementação real do contrato.
type taskService struct {
	// Depende da INTERFACE do Repository (não da struct concreta).
	// O Service não sabe se é Postgres, MySQL ou um banco falso de teste.
	repo repository.TaskRepository
}

// NewTaskService é o construtor do Service.
// Recebe o Repository já pronto e o guarda dentro do Service.
func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) List(ctx context.Context) ([]model.Task, error) {
	return s.repo.FindAll(ctx)
}

func (s *taskService) GetByID(ctx context.Context, id string) (model.Task, error) {
	task, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Task{}, ErrTaskNotFound
		}
		return model.Task{}, err
	}
	return task, nil
}

func (s *taskService) Create(ctx context.Context, input model.CreateTaskInput) (model.Task, error) {
	// strings.TrimSpace remove espaços do início e fim: "  abc  " → "abc"
	title := strings.TrimSpace(input.Title)
	description := strings.TrimSpace(input.Description)

	// Validações de negócio
	if title == "" {
		return model.Task{}, ErrTitleEmpty
	}
	if len([]rune(title)) < 3 {
		// len([]rune(...)) conta caracteres reais (funciona com acentos e emojis).
		// len(string) conta bytes — pode falhar com "ção" (3 letras, 5 bytes).
		return model.Task{}, ErrTitleTooShort
	}
	if description == "" {
		return model.Task{}, ErrDescEmpty
	}

	// Monta a tarefa com os campos que o SISTEMA gera (não o usuário).
	task := model.Task{
		ID:          uuid.New().String(), // UUID v4 aleatório
		Title:       title,
		Description: description,
		Done:        false,
		CreatedAt:   time.Now().Format(time.RFC3339), // Formato: 2024-01-15T10:30:00Z
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return model.Task{}, err
	}

	return task, nil
}

// Update busca a tarefa atual e aplica apenas os campos enviados.
func (s *taskService) Update(ctx context.Context, id string, input model.UpdateTaskInput) (model.Task, error) {
	// Reutiliza GetByID que já trata o "não encontrado"
	task, err := s.GetByID(ctx, id)
	if err != nil {
		return model.Task{}, err
	}

	// Só atualiza os campos que vieram na requisição.
	// O "*" antes de input.Title "desempacota" o ponteiro para ler o valor.
	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if len([]rune(title)) < 3 {
			return model.Task{}, ErrTitleTooShort
		}
		task.Title = title
	}

	if input.Description != nil {
		task.Description = strings.TrimSpace(*input.Description)
	}

	if err := s.repo.Update(ctx, task); err != nil {
		return model.Task{}, err
	}

	return task, nil
}

// Delete remove a tarefa se ela existir.
func (s *taskService) Delete(ctx context.Context, id string) error {
	found, err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if !found {
		return ErrTaskNotFound
	}
	return nil
}

// Count retorna o total de tarefas.
func (s *taskService) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}
