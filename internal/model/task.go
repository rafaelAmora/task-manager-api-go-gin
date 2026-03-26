package model

// As "tags" entre crases (``) são instruções para o Go:
//   - json:"title"  → quando virar JSON para o usuário, usa esse nome
//   - db:"title"    → quando buscar no banco, procura essa coluna

type Task struct {
	ID          string `db:"id"          json:"id"`
	Title       string `db:"title"       json:"title"`
	Description string `db:"description" json:"description"`
	Done        bool   `db:"done"        json:"done"`
	CreatedAt   string `db:"created_at"  json:"created_at"`
}

// CreateTaskInput são os dados que o usuário envia para CRIAR uma tarefa.
type CreateTaskInput struct {
	Title       string `json:"title"       binding:"required"`
	Description string `json:"description" binding:"required"`
}

// UpdateTaskInput são os dados que o usuário envia para ATUALIZAR uma tarefa.
//*string (com ponteiro) no Go é como uma variável no JS que pode ser null ou undefined.
type UpdateTaskInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
