package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rafaelAmora/db"
	"github.com/rafaelAmora/internal/controller"
	"github.com/rafaelAmora/internal/repository"
	"github.com/rafaelAmora/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Arquivo .env não encontrado — usando variáveis do sistema.")
	}

	database := db.Connect()

	// Cria o Repository — injetando a conexão com o banco
	taskRepo := repository.NewTaskRepository(database)

	//  Cria o Service — injetando o Repository
	taskSvc := service.NewTaskService(taskRepo)

	//  Cria o Controller — injetando o Service
	taskCtrl := controller.NewTaskController(taskSvc)

	//  Configura o roteador Gin
	router := gin.Default()

	//  Registra as rotas de tarefas
	taskCtrl.RegisterRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Servidor rodando em http://localhost:%s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
	}
}
