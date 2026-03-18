package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Task struct {
	Id          string `json:"id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Done        bool   `json:"done"`
	CreatedAt   string `json:"createdAt"`
}

var tasks = make(map[string]Task)

func GetTask(c *gin.Context) {

	var list []Task

	for _, t := range tasks {
		list = append(list, t)
	}

	c.JSON(http.StatusOK, list)
}

func CreateTask(c *gin.Context) {
	var newTask Task

	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTask.Id = uuid.NewString()

	newTask.CreatedAt = time.Now().Format(time.RFC3339)

	tasks[newTask.Id] = newTask

	c.JSON(http.StatusCreated, newTask)

}

func GetTaskById(c *gin.Context) {
	id := c.Param("id")
	if task, exists := tasks[id]; exists {
		c.JSON(http.StatusOK, task)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var newTask Task

	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if task, exists := tasks[id]; exists {
		task.Description = newTask.Description
		task.Title = newTask.Title
		task.Done = newTask.Done

		tasks[id] = task

		c.JSON(http.StatusOK, task)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

func Delete(c *gin.Context) {
	id := c.Param("id")

	if _, exists := tasks[id]; exists {
		delete(tasks, id)
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

func main() {
	router := gin.Default()
	router.GET("/tasks", GetTask)
	router.POST("/tasks", CreateTask)
	router.GET("/tasks/:id", GetTaskById)
	router.PATCH("/tasks/:id", UpdateTask)
	router.DELETE("/tasks/:id", Delete)

	router.Run("localhost:8080")
}
