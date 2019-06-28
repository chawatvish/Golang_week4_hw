package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func main() {
	r := gin.Default()
	p := fmt.Sprintf(":%s", os.Getenv("PORT"))

	r.GET("/api/todo", getTodosHandler)
	r.GET("/api/todos/:id", getTodoByID)
	r.POST("/api/todos", postTodoHandler)
	r.PUT("/api/todos/:id", updateTodoHandler)
	r.DELETE("/api/todos/:id", deleteTodoByID)

	r.Run(p)
}

func getTodosHandler(c *gin.Context) {
	db, err := connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	todos, err := queryTodos(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, todos)
}

func getTodoByID(c *gin.Context) {
	db, err := connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	todo, err := queryTodoByID(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, todo)
}

func postTodoHandler(c *gin.Context) {
	db, err := connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var todo Todo
	c.BindJSON(&todo)
	todo, err = addTodo(db, todo.Title, todo.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, todo)
}

func updateTodoHandler(c *gin.Context) {
	db, err := connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var todo Todo
	c.BindJSON(&todo)
	todo.ID = id
	err = updateTodo(db, id, todo.Title, todo.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, todo)
}

func deleteTodoByID(c *gin.Context) {
	db, err := connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = removeTodoByID(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}
