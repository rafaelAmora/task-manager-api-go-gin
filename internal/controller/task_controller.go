package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rafaelAmora/internal/model"
	"github.com/rafaelAmora/internal/service"
)

// Esta camada é a porta de entrada da API. Ela:
//   1. Recebe a requisição HTTP (método, URL, body, parâmetros)
//   2. Extrai os dados necessários
//   3. Chama o Service para fazer o trabalho
//   4. Devolve a resposta HTTP (status code + JSON)
//
// O Controller NÃO sabe nada de banco de dados.
// O Controller NÃO contém regras de negócio.
// Ele só "recebe", "delega" e "responde".

// TaskController guarda as dependências que o Controller precisa.
// No caso, só precisa do Service.
type TaskController struct {
	svc service.TaskService
}

// NewTaskController cria o Controller com o Service já injetado.
func NewTaskController(svc service.TaskService) *TaskController {
	return &TaskController{svc: svc}
}

func (ctrl *TaskController) RegisterRoutes(router *gin.Engine) {
	tasks := router.Group("/tasks")
	{
		tasks.GET("", ctrl.List)
		tasks.GET("/count", ctrl.Count) //Deve vir ANTES de /:id, senão o Gin confunde "count" com um id
		tasks.GET("/:id", ctrl.GetByID)
		tasks.POST("", ctrl.Create)
		tasks.PATCH("/:id", ctrl.Update)
		tasks.DELETE("/:id", ctrl.Delete)
	}
}

// handleError traduz erros de negócio em respostas HTTP.
// errors.Is() verifica se o erro é (ou contém) o erro esperado.
// É mais seguro que comparar strings de mensagem de erro.

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrTaskNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrTitleEmpty),
		errors.Is(err, service.ErrTitleTooShort),
		errors.Is(err, service.ErrDescEmpty):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno do servidor"})
	}
}

func (ctrl *TaskController) List(c *gin.Context) {
	tasks, err := ctrl.svc.List(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (ctrl *TaskController) GetByID(c *gin.Context) {
	id := c.Param("id") // c.Param("id") = req.params.id no Express

	task, err := ctrl.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

func (ctrl *TaskController) Create(c *gin.Context) {
	var input model.CreateTaskInput

	// ShouldBindJSON lê o body JSON e preenche a struct.
	// O `binding:"required"` nas tags da struct faz a validação básica.
	// No Express: req.body (com express.json() configurado)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "body inválido: " + err.Error()})
		return
	}

	task, err := ctrl.svc.Create(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, task) // 201 Created
}

// Atualiza parcialmente uma tarefa existente.
func (ctrl *TaskController) Update(c *gin.Context) {
	id := c.Param("id")

	var input model.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "body inválido"})
		return
	}

	task, err := ctrl.svc.Update(c.Request.Context(), id, input)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

func (ctrl *TaskController) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.svc.Delete(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}

	// 204 No Content = sucesso sem corpo na resposta
	c.JSON(http.StatusNoContent, nil)
}

// Controller → Service → Repository → Banco → Repository → Service → Controller
func (ctrl *TaskController) Count(c *gin.Context) {
	count, err := ctrl.svc.Count(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": count})
}
