package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgconn"
)

// Código SQLSTATE do PostgreSQL para violação de chave estrangeira.
const codigoViolacaoFK = "23503"

// mostrarErro imprime uma mensagem amigável quando o erro é a violação de uma
// regra de negócio garantida pelo banco; nos demais casos, imprime o erro original.
func mostrarErro(err error) {
	msg, ok := traduzirErro(err)
	if ok {
		fmt.Println(msg)
		return
	}

	fmt.Fprintln(os.Stderr, "erro:", err)
}

// traduzirErro procura um *pgconn.PgError dentro de err e, pelo nome da
// constraint violada, devolve a mensagem da regra de negócio correspondente.
func traduzirErro(err error) (string, bool) {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return "", false
	}

	switch pgErr.ConstraintName {
	case "ck_tipo_usuario":
		return "Tipo inválido (RN01): use porteiro, limpeza, supervisor ou morador.", true
	case "ck_feedback_emissor_receptor":
		return "Emissor e receptor devem ser diferentes (RN04).", true
	case "ck_status_feedback":
		return "Status inválido (RN14): use registrado, revisado ou arquivado.", true
	case "fk_feedback_emissor":
		return "Emissor não encontrado: informe o ID de um usuário cadastrado (RN03).", true
	case "fk_feedback_receptor":
		return "Receptor não encontrado: informe o ID de um usuário cadastrado (RN03).", true
	}

	return "", false
}

// usuarioTemFeedbacks indica se o erro veio de excluir um usuário que ainda é
// emissor ou receptor de algum feedback (RN15).
func usuarioTemFeedbacks(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == codigoViolacaoFK
}
