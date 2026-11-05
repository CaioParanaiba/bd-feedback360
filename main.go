package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	dbURL := os.Getenv("DATABASE_URL")

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro ao conectar ao PostgreSQL:", err)
		os.Exit(1)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	fmt.Println("Conectado!")
}
