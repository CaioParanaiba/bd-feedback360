package main

import (
	"database/sql"
	"fmt"
)

type usuario struct {
	ID           int64
	nome         string
	tipo_usuario string
}

func criarUsuario(db *sql.DB, nome, tipo string) error {
	_, err := db.Exec("INSERT INTO usuario (nome, tipo_usuario) VALUES ($1, $2)", nome, tipo)
	if err != nil {
		return fmt.Errorf("erro ao inserir usuario: %w", err)
	}

	return nil
}

func listarUsuarios(db *sql.DB) ([]usuario, error) {
	rows, err := db.Query("SELECT id_usuario, nome, tipo_usuario FROM usuario ORDER BY tasks.id")
	if err != nil {
		return nil, fmt.Errorf("erro na leitura de usuarios: %w", err)
	}

	defer rows.Close()

	var lista []usuario
	for rows.Next() {
		var v usuario
		err := rows.Scan(&v.ID, &v.nome, &v.tipo_usuario)
		if err != nil {
			return nil, fmt.Errorf("ler tarefa: %w", err)
		}
		lista = append(lista, v)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("percorrer tarefas: %w", err)
	}

	return lista, nil
}
