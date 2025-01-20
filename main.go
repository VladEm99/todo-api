package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL драйвер
)

// Task структура для задач
type Task struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Status      string  `json:"status"`
	Priority    int     `json:"priority"`
	DueDate     *string `json:"due_date"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func main() {
	// Подключение к базе данных
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		os.Getenv("HOST"), os.Getenv("DBPORT"), os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DBNAME"),
	))
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	// Проверка соединения
	if err = db.Ping(); err != nil {
		log.Fatalf("База данных недоступна: %v", err)
	}

	fmt.Println("Подключение к базе данных успешно установлено!")

	r := gin.Default()

	// GET /tasks - получение всех задач
	r.GET("/tasks", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name, description, status, priority, due_date, created_at, updated_at FROM tasks")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить задачи"})
			return
		}
		defer rows.Close()

		tasks := []Task{}
		for rows.Next() {
			var task Task
			var description, dueDate *string
			err = rows.Scan(&task.ID, &task.Name, &description, &task.Status, &task.Priority, &dueDate, &task.CreatedAt, &task.UpdatedAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения задач"})
				return
			}
			task.Description = description
			task.DueDate = dueDate
			tasks = append(tasks, task)
		}

		c.JSON(http.StatusOK, gin.H{"tasks": tasks})
	})

	// POST /tasks - добавление новой задачи
	r.POST("/tasks", func(c *gin.Context) {
		var newTask struct {
			Name        string  `json:"name"`
			Description *string `json:"description"`
			Status      string  `json:"status"`
			Priority    int     `json:"priority"`
			DueDate     *string `json:"due_date"`
		}
		if err := c.BindJSON(&newTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ввод данных"})
			return
		}

		query := `INSERT INTO tasks (id, name, description, status, priority, due_date)
				  VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)`
		_, err := db.Exec(query, newTask.Name, newTask.Description, newTask.Status, newTask.Priority, newTask.DueDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось добавить задачу"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Задача успешно добавлена"})
	})

	// Run server
	r.Run(":8080")
}