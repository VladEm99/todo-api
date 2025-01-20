package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	r := gin.Default()

	// In-memory storage for tasks
	tasks := map[string]string{}

	// Endpoint to get all tasks
	r.GET("/tasks", func(c *gin.Context) {
		c.JSON(200, gin.H{"tasks": tasks})
	})

	// Endpoint to add a new task
	r.POST("/tasks", func(c *gin.Context) {
		var newTask struct {
			Name string `json:"name"`
		}
		if err := c.BindJSON(&newTask); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}
		id := uuid.New().String()
		tasks[id] = newTask.Name
		c.JSON(201, gin.H{"id": id, "message": "Task added"})
	})

	// Run the server on port 8080
	r.Run(":8080")
}