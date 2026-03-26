# 📝 Task Manager API (Go + Gin)

![Go](https://img.shields.io/badge/Go-1.21-blue)
![Gin](https://img.shields.io/badge/Gin-Framework-lightgrey)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-blue)
![Status](https://img.shields.io/badge/status-em%20evolução-green)

Uma API REST de gerenciamento de tarefas desenvolvida em Go com o framework Gin e persistência de dados com PostgreSQL. O projeto evoluiu de um único arquivo `main.go` para uma **arquitetura em camadas** (Controller → Service → Repository), separando claramente as responsabilidades de cada parte do sistema.

---

## 🚀 Tecnologias Utilizadas

- [Go](https://golang.org/) — Linguagem principal
- [Gin](https://github.com/gin-gonic/gin) — Framework web para roteamento e tratamento de requisições HTTP
- [PostgreSQL](https://www.postgresql.org/) — Banco de dados relacional
- [sqlx](https://github.com/jmoiron/sqlx) — Extensão do `database/sql` para facilitar queries e mapeamento de structs
- [lib/pq](https://github.com/lib/pq) — Driver PostgreSQL para Go
- [Google UUID](https://github.com/google/uuid) — Geração de identificadores únicos
- [godotenv](https://github.com/joho/godotenv) — Carregamento de variáveis de ambiente via `.env`

---

## ✅ Funcionalidades

- [x] Criar uma nova tarefa
- [x] Listar todas as tarefas
- [x] Buscar uma tarefa pelo ID
- [x] Contar o total de tarefas
- [x] Atualizar parcialmente uma tarefa (PATCH)
- [x] Deletar uma tarefa

---

## 🗂️ Estrutura do Projeto

```
task-manager-api-go-gin/
├── .env
├── .gitignore
├── go.mod
├── go.sum
├── main.go
├── db/
│   └── postgres.go          # Conexão com o banco e criação da tabela
└── internal/
    ├── controller/
    │   └── task_controller.go   # Camada HTTP: recebe, delega e responde
    ├── model/
    │   └── task.go              # Structs de domínio e inputs
    ├── repository/
    │   └── task_repository.go   # Camada de dados: queries SQL
    └── service/
        └── task_service.go      # Regras de negócio e validações
```

---

## 🏗️ Arquitetura em Camadas

O projeto adota o padrão de separação de responsabilidades em 3 camadas, orquestradas pelo `main.go`:

```
HTTP Request
     │
     ▼
┌─────────────┐
│  Controller │  Recebe a requisição, extrai dados, devolve a resposta HTTP.
│             │  Não conhece o banco. Não contém regras de negócio.
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Service   │  Contém as regras de negócio e validações (ex: título mínimo de 3 chars).
│             │  Não conhece HTTP. Não escreve SQL.
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Repository  │  Única camada que fala com o banco.
│             │  Executa queries e mapeia resultados para structs.
└──────┬──────┘
       │
       ▼
  PostgreSQL
```

A comunicação entre as camadas é feita via **interfaces**, o que facilita a troca de implementações (ex: banco de dados diferente em testes).

---

## ▶️ Como Rodar o Projeto

### Pré-requisitos

- [Go 1.21+](https://golang.org/dl/) instalado
- [PostgreSQL](https://www.postgresql.org/download/) instalado e rodando

### Passo a passo

**1. Clone o repositório:**
```bash
git clone https://github.com/rafaelAmora/task-manager-api-go-gin.git
cd task-manager-api-go-gin
```

**2. Configure as variáveis de ambiente:**

Crie um arquivo `.env` na raiz do projeto:
```env
DB_NAME=seu_banco
DB_USER=seu_usuario
DB_PASSWORD=sua_senha
DB_HOST=localhost  # opcional, padrão: localhost
PORT=8080          # opcional, padrão: 8080
```

> A tabela `tasks` é criada automaticamente na primeira execução, caso não exista.

**3. Instale as dependências:**
```bash
go mod tidy
```

**4. Inicie o servidor:**
```bash
go run main.go
```

O servidor estará disponível em `http://localhost:8080`.

---

## 🛣️ Rotas da API

| Método   | Rota            | Descrição                          |
|----------|-----------------|------------------------------------|
| `GET`    | `/tasks`        | Lista todas as tarefas             |
| `GET`    | `/tasks/count`  | Retorna o total de tarefas         |
| `GET`    | `/tasks/:id`    | Busca uma tarefa pelo ID           |
| `POST`   | `/tasks`        | Cria uma nova tarefa               |
| `PATCH`  | `/tasks/:id`    | Atualiza parcialmente uma tarefa   |
| `DELETE` | `/tasks/:id`    | Deleta uma tarefa                  |

> ⚠️ A rota `/tasks/count` deve ser declarada **antes** de `/tasks/:id` para o Gin não interpretar `"count"` como um ID.

---

## 📦 Exemplos de Uso

### Criar uma tarefa — `POST /tasks`

**Request body:**
```json
{
  "title": "Estudar Go",
  "description": "Aprender sobre structs, interfaces e goroutines"
}
```

**Response `201 Created`:**
```json
{
  "id": "a3f1c2d4-...",
  "title": "Estudar Go",
  "description": "Aprender sobre structs, interfaces e goroutines",
  "done": false,
  "created_at": "2025-01-15T10:30:00Z"
}
```

---

### Listar todas as tarefas — `GET /tasks`

**Response `200 OK`:**
```json
[
  {
    "id": "a3f1c2d4-...",
    "title": "Estudar Go",
    "description": "Aprender sobre structs, interfaces e goroutines",
    "done": false,
    "created_at": "2025-01-15T10:30:00Z"
  }
]
```

---

### Contar tarefas — `GET /tasks/count`

**Response `200 OK`:**
```json
{
  "total": 5
}
```

---

### Buscar tarefa por ID — `GET /tasks/:id`

**Response `200 OK`:**
```json
{
  "id": "a3f1c2d4-...",
  "title": "Estudar Go",
  "description": "Aprender sobre structs, interfaces e goroutines",
  "done": false,
  "created_at": "2025-01-15T10:30:00Z"
}
```

**Response `404 Not Found`:**
```json
{
  "error": "tarefa não encontrada"
}
```

---

### Atualizar uma tarefa — `PATCH /tasks/:id`

Apenas os campos enviados são atualizados.

**Request body:**
```json
{
  "description": "Revisar goroutines e channels"
}
```

**Response `200 OK`:**
```json
{
  "id": "a3f1c2d4-...",
  "title": "Estudar Go",
  "description": "Revisar goroutines e channels",
  "done": false,
  "created_at": "2025-01-15T10:30:00Z"
}
```

---

### Deletar uma tarefa — `DELETE /tasks/:id`

**Response `204 No Content`**

---

## ⚠️ Validações e Erros

O Service aplica validações de negócio antes de salvar os dados:

| Situação                        | Status HTTP | Mensagem de erro                          |
|---------------------------------|-------------|-------------------------------------------|
| Título vazio                    | `422`       | `o título não pode estar vazio`           |
| Título com menos de 3 caracteres | `422`      | `o título precisa ter pelo menos 3 caracteres` |
| Descrição vazia                 | `422`       | `a descrição não pode estar vazia`        |
| Tarefa não encontrada           | `404`       | `tarefa não encontrada`                   |
| Erro interno                    | `500`       | `erro interno do servidor`                |

---

## 🔮 Melhorias Futuras

- [ ] **Autenticação** — Implementar JWT para proteger as rotas
- [ ] **Testes** — Adicionar testes unitários e de integração (as interfaces já facilitam o mock do Repository)
- [ ] **Dockerização** — Containerizar a aplicação com Docker e Docker Compose (incluindo o PostgreSQL)
- [ ] **Paginação** — Suporte a paginação na listagem de tarefas
- [ ] **Migrations** — Gerenciar o schema do banco com uma ferramenta de migrations (ex: `golang-migrate`)

---

## 💬 Feedback

Este é um projeto de aprendizado em constante evolução. Sugestões de boas práticas e críticas construtivas são muito bem-vindas!

Abra uma [issue](https://github.com/rafaelAmora/task-manager-api-go-gin/issues) ou entre em contato. 🚀