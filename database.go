package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func connect() (*sql.DB, error) {
	return sql.Open("postgres", os.Getenv("DATABASE_URL"))
}

func queryTodos(db *sql.DB) ([]Todo, error) {
	stmt, err := db.Prepare("SELECT id, title, status FROM todos")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Status)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func queryTodoByID(db *sql.DB, id int) (Todo, error) {
	stmt, err := db.Prepare("SELECT id, title, status FROM todos WHERE id=$1;")
	if err != nil {
		return Todo{}, err
	}
	rows := stmt.QueryRow(id)

	var todo Todo
	err = rows.Scan(&todo.ID, &todo.Title, &todo.Status)
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func addTodo(db *sql.DB, title string, status string) (int, error) {
	query := `
	INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING ID
	`
	var id int
	row := db.QueryRow(query, title, status)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func updateTodoStatus(db *sql.DB, id int, status string) error {
	stmt, err := db.Prepare("UPDATE todos SET status=$2 WHERE id=$1;")
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(id, status); err != nil {
		return err
	}

	return nil
}

func updateTodoTitle(db *sql.DB, id int, title string) error {
	stmt, err := db.Prepare("UPDATE todos SET title=$2 WHERE id=$1;")
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(id, title); err != nil {
		return err
	}

	return nil
}

func removeTodoByID(db *sql.DB, id int) error {
	stmt, err := db.Prepare("DELETE FROM todos WHERE id=$1;")
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(id); err != nil {
		return err
	}

	return nil
}
