package main

import (
	"database/sql"
	"fmt"
	"time"
)

type Feedback struct {
	ID         int64
	Descricao  string
	DataHora   time.Time
	Status     string
	IDEmissor  int64
	IDReceptor int64
}

func criarFeedback(db *sql.DB, feed Feedback) (int64, error) {
	var id int64
	err := db.QueryRow("INSERT INTO feedback (descricao, id_emissor, id_receptor) VALUES ($1,$2,$3) RETURNING id_feedback",
		feed.Descricao, feed.IDEmissor, feed.IDReceptor).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("inserir feedback: %w", err)
	}

	return id, nil
}

func buscarFeedback(db *sql.DB, id int64) (Feedback, error) {
	var feed Feedback
	err := db.QueryRow("SELECT id_feedback, descricao, data_hora, status, id_emissor, id_receptor FROM feedback WHERE id_feedback = $1", id).Scan(
		&feed.ID, &feed.Descricao, &feed.DataHora, &feed.Status, &feed.IDEmissor, &feed.IDReceptor)
	if err != nil {
		return Feedback{}, fmt.Errorf("buscar feedback %d: %w", id, err)
	}

	return feed, nil
}

func editarFeedback(db *sql.DB, feed Feedback) error {
	_, err := db.Exec(`UPDATE feedback SET descricao = $1, status = $2, id_emissor = $3, id_receptor = $4
						WHERE id_feedback = $5`, feed.Descricao, feed.Status, feed.IDEmissor, feed.IDReceptor, feed.ID)
	if err != nil {
		return fmt.Errorf("atualizar feedback %d: %w", feed.ID, err)
	}

	return nil
}

func excluirFeedback(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM feedback WHERE id_feedback = $1", id)
	if err != nil {
		return fmt.Errorf("deletar id %d: %w", id, err)
	}

	return nil
}

func listarFeedbacksPorReceptor(db *sql.DB, idReceptor int64) ([]Feedback, error) {
	rows, err := db.Query(`SELECT id_feedback, descricao, data_hora, status, id_emissor, id_receptor
		FROM feedback WHERE id_receptor = $1 ORDER BY data_hora`, idReceptor)
	if err != nil {
		return nil, fmt.Errorf("leitura de feedbacks: %w", err)
	}

	defer rows.Close()

	var lista []Feedback
	for rows.Next() {
		var f Feedback
		err := rows.Scan(&f.ID, &f.Descricao, &f.DataHora, &f.Status, &f.IDEmissor, &f.IDReceptor)
		if err != nil {
			return nil, fmt.Errorf("ler feedback: %w", err)
		}
		lista = append(lista, f)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("percorrer feedbacks: %w", err)
	}

	return lista, nil
}
