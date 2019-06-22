package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func connect() (*sql.DB, error) {
	return sql.Open("postgres", os.Getenv("DATABASE_URL"))
}

func queryTodos() ([]Todo, error) {
	db, err := connect()
	if err != nil {
		db.Close()
		return nil, err
	}

	stmt, err := db.Prepare("SELECT id, title, status FROM todos")
	if err != nil {
		db.Close()
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		db.Close()
		return nil, err
	}

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Status)
		if err != nil {
			db.Close()
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func queryTodoByID(id int) (Todo, error) {
	db, err := connect()
	if err != nil {
		db.Close()
		return Todo{}, err
	}

	stmt, err := db.Prepare("SELECT id, title, status FROM todos")
	if err != nil {
		db.Close()
		return Todo{}, err
	}
	rows, err := stmt.Query()
	if err != nil {
		db.Close()
		return Todo{}, err
	}

	var todo Todo
	err = rows.Scan(&todo.ID, &todo.Title, &todo.Status)
	if err != nil {
		db.Close()
		return Todo{}, err
	}

	return todo, nil
}

func addTodo(title string, status string) (Todo, error) {
	db, err := connect()
	if err != nil {
		db.Close()
		return Todo{}, err
	}

	query := `
	INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING ID
	`
	var newTodo Todo
	row := db.QueryRow(query, title, status)
	err = row.Scan(&newTodo.ID, &newTodo.Title, &newTodo.Status)
	if err != nil {
		return Todo{}, err
	}

	return newTodo, nil
}

func updateTodoStatus(id int, status string) error {
	db, err := connect()
	if err != nil {
		db.Close()
		return err
	}

	stmt, err := db.Prepare("UPDATE todos SET status=$2 WHERE id=$1;")
	if err != nil {
		db.Close()
		return err
	}

	if _, err := stmt.Exec(id, status); err != nil {
		db.Close()
		return err
	}

	return nil
}

func updateTodoTitle(id int, title string) error {
	db, err := connect()
	if err != nil {
		db.Close()
		return err
	}

	stmt, err := db.Prepare("UPDATE todos SET title=$2 WHERE id=$1;")
	if err != nil {
		db.Close()
		return err
	}

	if _, err := stmt.Exec(id, title); err != nil {
		db.Close()
		return err
	}

	return nil
}
