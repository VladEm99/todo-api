package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Task struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := gin.Default()

	// In-memory storage for tasks
	var tasks []Task

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

		// Создаем новую задачу с уникальным ID
		task := Task{
			ID:   uuid.New().String(),
			Name: newTask.Name,
		}
		tasks = append(tasks, task)
		c.JSON(201, gin.H{"message": "Task added", "task": task})
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

		// Поиск задачи по ID
		for i, task := range tasks {
			if task.ID == id {
				tasks[i].Name = updatedTask.Name
				c.JSON(200, gin.H{"message": "Task updated", "task": tasks[i]})
				return
			}
		}

		c.JSON(404, gin.H{"error": "Task not found"})
	})

	// DELETE /tasks/:id - удаляет задачу
	r.DELETE("/tasks/:id", func(c *gin.Context) {
		id := c.Param("id") // Получаем id из параметра пути

		// Поиск задачи по ID и удаление
		for i, task := range tasks {
			if task.ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				c.JSON(200, gin.H{"message": "Task deleted"})
				return
			}
		}

		c.JSON(404, gin.H{"error": "Task not found"})
	})

	// Run the server on port 8080
	r.Run(":8080")
}