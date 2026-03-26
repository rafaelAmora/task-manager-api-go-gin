package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // O "_" importa só para registrar o driver do Postgres.
	// Sem esse import, o sqlx não sabe como falar com o Postgres.
)

func Connect() *sqlx.DB {
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")

	if dbHost == "" {
		dbHost = "localhost" // valor padrão se não configurado
	}

	// DSN = Data Source Name. É a "string de endereço" do banco.

	dsn := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName)

	db, err := sqlx.Connect("postgres", dsn)

	if err != nil {
		log.Fatalf("❌ Erro ao conectar no banco: %v", err)
	}

	// Cria a tabela se ela ainda não existir.
	schema := `
	CREATE TABLE IF NOT EXISTS tasks (
		id          TEXT PRIMARY KEY,
		title       TEXT NOT NULL,
		description TEXT NOT NULL,
		done        BOOLEAN NOT NULL DEFAULT FALSE,
		created_at  TEXT NOT NULL
	);`

	db.MustExec(schema) // MustExec = executar consultas SQL que não retornam linhas
	log.Println("✅ Banco conectado e tabela verificada.")

	return db
}
