package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

const menus = `
╔══════════════════════════════════════════╗
║          FEEDBACK 360° — CONDOMÍNIO      ║
╠══════════════════════════════════════════╣
║  USUÁRIOS              FEEDBACKS         ║
║  1 · Cadastrar         5 · Cadastrar     ║
║  2 · Listar            6 · Por receptor  ║
║  3 · Editar            7 · Editar        ║
║  4 · Excluir           8 · Excluir       ║
╠══════════════════════════════════════════╣
║  0 · Sair                                ║
╚══════════════════════════════════════════╝
Opção: `

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

	menu(db)

}

func menu(db *sql.DB) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(menus)

		if !scanner.Scan() {
			return
		}

		opcao := strings.TrimSpace(scanner.Text())

		switch opcao {
		case "0":
			return

		case "1":
			fmt.Print("Nome: ")
			scanner.Scan()
			nome := strings.TrimSpace(scanner.Text())

			fmt.Print("Tipo (porteiro, limpeza, supervisor, morador): ")
			scanner.Scan()
			tipo := strings.TrimSpace(scanner.Text())

			id, err := criarUsuario(db, nome, tipo)
			if err != nil {
				fmt.Fprintln(os.Stderr, "erro:", err)
				continue
			}

			fmt.Println("Usuário cadastrado com ID", id)

		case "2":
			lista, err := listarUsuarios(db)
			if err != nil {
				fmt.Fprintln(os.Stderr, "erro:", err)
				continue
			}

			if len(lista) == 0 {
				fmt.Println("Nenhum usuário cadastrado.")
				continue
			}

			fmt.Printf("%-4s %-30s %s\n", "ID", "NOME", "TIPO")
			for _, u := range lista {
				fmt.Printf("%-4d %-30s %s\n", u.ID, u.Nome, u.TipoUsuario)
			}

		default:
			fmt.Println("Opcao inexistente")
		}
	}
}
