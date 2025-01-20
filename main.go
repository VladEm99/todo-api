package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/gin-gonic/gin"
)

var db *pgx.Conn

func main() {
	// Подключение к базе данных
	databaseURL := os.Getenv("DATABASE_URL") // Читаем URL базы из переменной окружения
	var err error
	db, err = pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close(context.Background())
	log.Println("Connected to the database successfully!")

	// Настройка Gin
	r := gin.Default()

	// Endpoints
	r.GET("/tasks", getTasks)
	r.POST("/tasks", createTask)

	// Запуск сервера
	r.Run(":8080")
}

// Пример: Получить все задачи
func getTasks(c *gin.Context) {
	rows, err := db.Query(context.Background(), "SELECT id, name FROM tasks")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to query tasks"})
		return
	}
	defer rows.Close()

	var tasks []map[string]interface{}
	for rows.Next() {
		var id string
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			c.JSON(500, gin.H{"error": "Failed to scan task"})
			return
		}
		tasks = append(tasks, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}

	c.JSON(200, gin.H{"tasks": tasks})
}

// Пример: Создать новую задачу
func createTask(c *gin.Context) {
	var newTask struct {
		Name string `json:"name"`
	}
	if err := c.BindJSON(&newTask); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	_, err := db.Exec(context.Background(), "INSERT INTO tasks (name) VALUES ($1)", newTask.Name)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to insert task"})
		return
	}

	c.JSON(201, gin.H{"message": "Task created"})
}