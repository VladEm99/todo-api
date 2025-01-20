package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Task represents a task with an ID and Name
type Task struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := gin.Default()

	// In-memory storage for tasks
	tasks := []Task{}

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

    	// Генерируем UUID
    	taskID := uuid.New().String()
    	task := Task{ID: taskID, Name: newTask.Name}

    	// Добавляем задачу в массив
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

		// Найти задачу по ID
		for i, task := range tasks {
			if task.ID == id {
				// Обновляем задачу
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

		// Найти задачу по ID
		for i, task := range tasks {
			if task.ID == id {
				// Удаляем задачу
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