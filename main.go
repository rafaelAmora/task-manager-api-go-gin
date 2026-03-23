package main

import (
	"fmt"
	"net/http"

	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type Task struct {
	Id          string `db:"id" json:"id"`
	Title       string `db:"title" json:"title" binding:"required"`
	Description string `db:"description" json:"description" binding:"required"`
	Done        bool   `db:"done" json:"done"`
	CreatedAt   string `db:"created_at" json:"created_at"`
}

// json:"createdAt": Diz ao Gin: "Quando você for transformar essa struct em um JSON para enviar ao usuário (ou receber dele), use o nome createdAt".

// db:"created_at": Diz ao sqlx: "Quando você for buscar dados no banco ou salvar lá, procure por uma coluna chamada created_at".

var DB *sqlx.DB

func GetTask(c *gin.Context) {

	var tasks []Task

	err := DB.Select(&tasks, `SELECT * FROM tasks`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {

	var newTask Task

	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTask.Id = uuid.NewString()

	newTask.CreatedAt = time.Now().Format(time.RFC3339)

	query := `INSERT INTO tasks (id, title, description, done, created_at)
	          VALUES (:id, :title, :description, :done, :created_at)`

	_, err := DB.NamedExec(query, newTask)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newTask)

}

func GetTaskById(c *gin.Context) {
	id := c.Param("id")

	query := `SELECT * FROM tasks WHERE id = $1 `

	var task Task

	err := DB.Get(&task, query, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)

}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var updateTask Task

	if err := c.ShouldBindJSON(&updateTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE tasks SET title = $1, description = $2 WHERE id  = $3 `

	result, err := DB.Exec(query, updateTask.Title, updateTask.Description, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// RowsAffected verifica se alguma linha foi de fato atualizada
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Se nenhuma linha foi afetada, o id não existe no banco
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

func Delete(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM tasks WHERE id = $1`

	result, err := DB.Exec(query, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar arquivo .env")
	}

	ConnectDB()

	router := gin.Default()
	router.GET("/tasks", GetTask)
	router.POST("/tasks", CreateTask)
	router.GET("/tasks/:id", GetTaskById)
	router.PATCH("/tasks/:id", UpdateTask)
	router.DELETE("/tasks/:id", Delete)

	router.Run("localhost:8080")
}

func ConnectDB() {
	var err error

	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("host=localhost port=5432 user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)

	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}
}
