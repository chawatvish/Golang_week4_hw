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

	r.GET("/api/todos", getTodosHandler)
	r.GET("/api/todos/:id", getTodoByID)
	r.POST("/api/todos", postTodoHandler)
	r.DELETE("/api/todos/:id", deleteTodoByID)

	r.Run(p)
}

func getTodosHandler(c *gin.Context) {
	todos, err := queryTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, todos)
}

func getTodoByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	todo, err := queryTodoByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, todo)
}

func postTodoHandler(c *gin.Context) {
	var todo Todo
	c.BindJSON(&todo)
	id, err := addTodo(todo.Title, todo.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, gin.H{"status": fmt.Sprintf("ID %d Added", id)})
}

func deleteTodoByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if err := removeTodoByID(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, gin.H{"status": "Deleted"})
}
