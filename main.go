package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID     int
	Title  string
	Status string
}

func main() {
	r := gin.Default()

	r.GET("/student", getStudentHandler)
	r.GET("/students", postStudentHandler)

	r.GET("/api/todos", getTodos)
	r.Run(":1234")
}

func getStudentHandler(c *gin.Context) {

}

func postStudentHandler(c *gin.Context) {

}

func getTodos(c *gin.Context) {
	todos, err := queryTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, todos)
}
