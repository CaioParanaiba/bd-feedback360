package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
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
				mostrarErro(err)
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
		case "3":
			fmt.Print("ID do usuario: ")
			scanner.Scan()
			idTexto := strings.TrimSpace(scanner.Text())

			id, err := strconv.ParseInt(idTexto, 10, 64)
			if err != nil {
				fmt.Println("ID invalido")
				continue
			}

			u, err := buscarUsuario(db, id)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					fmt.Println("Usuário não encontrado")
				} else {
					fmt.Fprintln(os.Stderr, "erro:", err)
				}
				continue
			}
			fmt.Printf("Nome [%s]: ", u.Nome)
			scanner.Scan()
			nome := strings.TrimSpace(scanner.Text())
			if nome != "" {
				u.Nome = nome
			}

			fmt.Printf("Tipo [%s]: ", u.TipoUsuario)
			scanner.Scan()
			tipo := strings.TrimSpace(scanner.Text())
			if tipo != "" {
				u.TipoUsuario = tipo
			}

			err = editarUsuario(db, u)
			if err != nil {
				mostrarErro(err)
				continue
			}

			fmt.Println("Usuário atualizado")
		case "4":
			fmt.Print("ID do usuario: ")
			scanner.Scan()
			idTexto := strings.TrimSpace(scanner.Text())

			id, err := strconv.ParseInt(idTexto, 10, 64)
			if err != nil {
				fmt.Println("ID invalido")
				continue
			}

			u, err := buscarUsuario(db, id)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					fmt.Println("Usuário não encontrado")
				} else {
					fmt.Fprintln(os.Stderr, "erro:", err)
				}
				continue
			}

			err = excluirUsuario(db, u.ID)
			if err != nil {
				if usuarioTemFeedbacks(err) {
					fmt.Println("Usuário possui feedbacks e não pode ser excluído (RN15).")
				} else {
					fmt.Fprintln(os.Stderr, "erro:", err)
				}
				continue
			}
			fmt.Println("Usuário excluido")

		case "5":
			fmt.Print("Descrição: ")
			scanner.Scan()
			descricao := strings.TrimSpace(scanner.Text())

			fmt.Print("ID do emissor: ")
			scanner.Scan()
			idEmissor, err := strconv.ParseInt(strings.TrimSpace(scanner.Text()), 10, 64)
			if err != nil {
				fmt.Println("ID invalido")
				continue
			}

			fmt.Print("ID do receptor: ")
			scanner.Scan()
			idReceptor, err := strconv.ParseInt(strings.TrimSpace(scanner.Text()), 10, 64)
			if err != nil {
				fmt.Println("ID invalido")
				continue
			}

			feed := Feedback{Descricao: descricao, IDEmissor: idEmissor, IDReceptor: idReceptor}
			id, err := criarFeedback(db, feed)
			if err != nil {
				mostrarErro(err)
				continue
			}

			fmt.Println("Feedback cadastrado com ID", id)

		case "6":
			fmt.Print("ID do receptor: ")
			scanner.Scan()
			idReceptor, err := strconv.ParseInt(strings.TrimSpace(scanner.Text()), 10, 64)
			if err != nil {
				fmt.Println("ID invalido")
				continue
			}

			lista, err := listarFeedbacksPorReceptor(db, idReceptor)
			if err != nil {
				fmt.Fprintln(os.Stderr, "erro:", err)
				continue
			}

			if len(lista) == 0 {
				fmt.Println("Nenhum feedback para esse receptor.")
				continue
			}

			fmt.Printf("%-4s %-16s %-11s %-8s %s\n", "ID", "DATA", "STATUS", "EMISSOR", "DESCRIÇÃO")
			for _, f := range lista {
				fmt.Printf("%-4d %-16s %-11s %-8d %s\n", f.ID, f.DataHora.Format("02/01/2006 15:04"), f.Status, f.IDEmissor, f.Descricao)
			}

		case "7":
			fmt.Print("ID do feedback: ")
			scanner.Scan()
			id, err := strconv.ParseInt(strings.TrimSpace(scanner.Text()), 10, 64)
			if err != nil {
				fmt.Println("ID invalido")
				continue
			}

			feed, err := buscarFeedback(db, id)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					fmt.Println("Feedback não encontrado")
				} else {
					fmt.Fprintln(os.Stderr, "erro:", err)
				}
				continue
			}

			fmt.Printf("Descrição [%s]: ", feed.Descricao)
			scanner.Scan()
			descricao := strings.TrimSpace(scanner.Text())
			if descricao != "" {
				feed.Descricao = descricao
			}

			fmt.Printf("Status (registrado, revisado, arquivado) [%s]: ", feed.Status)
			scanner.Scan()
			status := strings.TrimSpace(scanner.Text())
			if status != "" {
				feed.Status = status
			}

			fmt.Printf("ID do emissor [%d]: ", feed.IDEmissor)
			scanner.Scan()
			texto := strings.TrimSpace(scanner.Text())
			if texto != "" {
				idEmissor, err := strconv.ParseInt(texto, 10, 64)
				if err != nil {
					fmt.Println("ID invalido")
					continue
				}
				feed.IDEmissor = idEmissor
			}

			fmt.Printf("ID do receptor [%d]: ", feed.IDReceptor)
			scanner.Scan()
			texto = strings.TrimSpace(scanner.Text())
			if texto != "" {
				idReceptor, err := strconv.ParseInt(texto, 10, 64)
				if err != nil {
					fmt.Println("ID invalido")
					continue
				}
				feed.IDReceptor = idReceptor
			}

			err = editarFeedback(db, feed)
			if err != nil {
				mostrarErro(err)
				continue
			}

			fmt.Println("Feedback atualizado")

		case "8":
			fmt.Print("ID do feedback: ")
			scanner.Scan()
			id, err := strconv.ParseInt(strings.TrimSpace(scanner.Text()), 10, 64)
			if err != nil {
				fmt.Println("ID invalido")
				continue
			}

			feed, err := buscarFeedback(db, id)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					fmt.Println("Feedback não encontrado")
				} else {
					fmt.Fprintln(os.Stderr, "erro:", err)
				}
				continue
			}

			err = excluirFeedback(db, feed.ID)
			if err != nil {
				fmt.Fprintln(os.Stderr, "erro:", err)
				continue
			}
			fmt.Println("Feedback excluido")

		default:
			fmt.Println("Opcao inexistente")
		}
	}
}
