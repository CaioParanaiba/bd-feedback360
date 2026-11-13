package main

import (
	"database/sql"
	"fmt"
)

type Usuario struct {
	ID          int64
	Nome        string
	TipoUsuario string
}

func criarUsuario(db *sql.DB, nome, tipo string) (int64, error) {
	var id int64
	err := db.QueryRow("INSERT INTO usuario (nome, tipo_usuario) VALUES ($1, $2) RETURNING id_usuario", nome, tipo).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("inserir usuario: %w", err)
	}

	return id, nil
}

func listarUsuarios(db *sql.DB) ([]Usuario, error) {
	rows, err := db.Query("SELECT id_usuario, nome, tipo_usuario FROM usuario ORDER BY id_usuario")
	if err != nil {
		return nil, fmt.Errorf("leitura de usuarios: %w", err)
	}

	defer rows.Close()

	var lista []Usuario
	for rows.Next() {
		var v Usuario
		err := rows.Scan(&v.ID, &v.Nome, &v.TipoUsuario)
		if err != nil {
			return nil, fmt.Errorf("ler usuario: %w", err)
		}
		lista = append(lista, v)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("percorrer usuarios: %w", err)
	}

	return lista, nil
}

func buscarUsuario(db *sql.DB, id int64) (Usuario, error) {
	var u Usuario
	err := db.QueryRow("SELECT id_usuario, nome, tipo_usuario FROM usuario WHERE id_usuario = $1", id).Scan(&u.ID, &u.Nome, &u.TipoUsuario)
	if err != nil {
		return Usuario{}, fmt.Errorf("buscar usuario %d: %w", id, err)
	}

	return u, nil
}

func editarUsuario(db *sql.DB, u Usuario) error {
	_, err := db.Exec("UPDATE usuario SET nome = $1, tipo_usuario = $2 WHERE id_usuario = $3", u.Nome, u.TipoUsuario, u.ID)
	if err != nil {
		return fmt.Errorf("atualizar usuario %d: %w", u.ID, err)
	}

	return nil
}

func excluirUsuario(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM usuario WHERE id_usuario = $1", id)
	if err != nil {
		return fmt.Errorf("deletar id %d: %w", id, err)
	}

	return nil
}
