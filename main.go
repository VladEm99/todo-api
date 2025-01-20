package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// In-memory storage for tasks
	tasks := []string{}

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
		tasks = append(tasks, newTask.Name)
		c.JSON(201, gin.H{"message": "Task added"})
	})

	// PUT /tasks/:id - обновляет задачу
	r.PUT("/tasks/:id", func(c *gin.Context) {
		id := c.Param("id") // Получаем id из параметра пути
		var updatedTask struct {
			Name string `json:"name"`
		}
		if err := c.BindJSON(&updatedTask); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		// Преобразуем id в индекс
		taskIndex, err := strconv.Atoi(id)
		if err != nil || taskIndex < 0 || taskIndex >= len(tasks) {
			c.JSON(404, gin.H{"error": "Task not found"})
			return
		}

		// Обновляем задачу
		tasks[taskIndex] = updatedTask.Name
		c.JSON(200, gin.H{"message": "Task updated"})
	})

	// DELETE /tasks/:id - удаляет задачу
	r.DELETE("/tasks/:id", func(c *gin.Context) {
		id := c.Param("id") // Получаем id из параметра пути

		// Преобразуем id в индекс
		taskIndex, err := strconv.Atoi(id)
		if err != nil || taskIndex < 0 || taskIndex >= len(tasks) {
			c.JSON(404, gin.H{"error": "Task not found"})
			return
		}

		// Удаляем задачу
		tasks = append(tasks[:taskIndex], tasks[taskIndex+1:]...)
		c.JSON(200, gin.H{"message": "Task deleted"})
	})

	// Run the server on port 8080
	r.Run(":8080")
}
