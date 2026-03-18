# 📝 Task Manager API (Go + Gin)

![Go](https://img.shields.io/badge/Go-1.21-blue)
![Gin](https://img.shields.io/badge/Gin-Framework-lightgrey)
![Status](https://img.shields.io/badge/status-learning-green)

Uma API REST simples de gerenciamento de tarefas, desenvolvida em Go com o framework Gin. Este projeto foi criado com fins de aprendizado, explorando conceitos de desenvolvimento backend como roteamento HTTP, manipulação de JSON e estruturação de uma API RESTful.

---

## 🚀 Tecnologias Utilizadas

- [Go](https://golang.org/) — Linguagem principal
- [Gin](https://github.com/gin-gonic/gin) — Framework web para roteamento e tratamento de requisições HTTP
- [Google UUID](https://github.com/google/uuid) — Geração de identificadores únicos
- Armazenamento em memória (`map`) — Sem banco de dados por enquanto

---

## ✅ Funcionalidades

- [x] Criar uma nova tarefa
- [x] Listar todas as tarefas
- [x] Buscar uma tarefa pelo ID
- [x] Atualizar os dados de uma tarefa
- [x] Deletar uma tarefa

---

## 🗂️ Estrutura do Projeto

```
taskcrud/
├── go.mod
├── go.sum
└── main.go
```

> Por ser um projeto introdutório, toda a lógica está concentrada em `main.go`. A separação em camadas (handlers, services, repositories) está listada como melhoria futura.

---

## ▶️ Como Rodar o Projeto

### Pré-requisitos

- [Go 1.21+](https://golang.org/dl/) instalado na máquina

### Passo a passo

**1. Clone o repositório:**
```bash
git clone https://github.com/rafaelAmora/taskcrud.git
cd taskcrud
```

**2. Instale as dependências:**
```bash
go mod tidy
```

**3. Inicie o servidor:**
```bash
go run main.go
```

O servidor estará disponível em `http://localhost:8080`.

---

## 🛣️ Rotas da API

| Método   | Rota            | Descrição                     |
|----------|-----------------|-------------------------------|
| `GET`    | `/tasks`        | Lista todas as tarefas        |
| `POST`   | `/tasks`        | Cria uma nova tarefa          |
| `GET`    | `/tasks/:id`    | Busca uma tarefa pelo ID      |
| `PATCH`  | `/tasks/:id`    | Atualiza uma tarefa existente |
| `DELETE` | `/tasks/:id`    | Deleta uma tarefa             |

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
  "createdAt": "2025-01-15T10:30:00Z"
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
    "createdAt": "2025-01-15T10:30:00Z"
  }
]
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
  "createdAt": "2025-01-15T10:30:00Z"
}
```

**Response `404 Not Found`:**
```json
{
  "message": "task not found"
}
```

---

### Atualizar uma tarefa — `PATCH /tasks/:id`

**Request body:**
```json
{
  "title": "Estudar Go",
  "description": "Revisar goroutines e channels",
  "done": true
}
```

**Response `200 OK`:**
```json
{
  "id": "a3f1c2d4-...",
  "title": "Estudar Go",
  "description": "Revisar goroutines e channels",
  "done": true,
  "createdAt": "2025-01-15T10:30:00Z"
}
```

---

### Deletar uma tarefa — `DELETE /tasks/:id`

**Response `204 No Content`** (sem corpo na resposta)

---

## 🔮 Melhorias Futuras

À medida que for evoluindo nos estudos, pretendo incorporar as seguintes melhorias:

- [ ] **Banco de dados** — Integração com PostgreSQL ou SQLite usando GORM
- [ ] **Arquitetura em camadas** — Separar o projeto em `handlers`, `services` e `repositories`
- [ ] **Autenticação** — Implementar JWT para proteger as rotas
- [ ] **Variáveis de ambiente** — Usar `.env` para configurações sensíveis
- [ ] **Testes** — Adicionar testes unitários e de integração
- [ ] **Dockerização** — Containerizar a aplicação com Docker
- [ ] **Paginação** — Suporte a paginação na listagem de tarefas

---

## 💬 Feedback

Este é um projeto de aprendizado, então estou aberto a sugestões de melhoria, boas práticas ou críticas construtivas.

Se quiser contribuir ou apontar algo que pode ser melhorado, fique à vontade para abrir uma [issue](https://github.com/rafaelAmora/taskcrud/issues) ou entrar em contato.
