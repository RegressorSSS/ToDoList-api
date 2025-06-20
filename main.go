package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/RegressorSSS/todolist/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

func main() {

	dbURL := "postgres://postgres:postgres@localhost:5432/postgres"
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %s", err)
		return
	}

	err = conn.Ping(conext.Background())
	if err != nil {
		log.Fatalf("ошибка при пинге БД: %s", err)
	}

	tasksHandler := handler.New(conn)

	r := gin.Default()

	// создать новую задачу
	r.POST("/tasks", HandleCreateTask)
	// получить все задачи
	r.GET("/tasks", HandleGetAllTasks)
	// поменять значение выполненности
	r.POST("/tasks/:taskId", HandleToggleComplete)
	// получить 1 задачу по ее id
	r.GET("/tasks/:taskId", HandleGetTask)
	// удалить задачу по ее id
	r.DELETE("/tasks/:taskId", HandleDeleteTask)
	// изменить задачу по ее id
	r.PATCH("/tasks/:taskId", HandleChangeTask)

	r.Run()

}

func HandleCreateTask(c *gin.Context) {
	var task Task
	err := json.NewDecoder(c.Request.Body).Decode(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У тебя неправильный json"+err.Error())
		return
	}

	id := uuid.New()
	task.Id = id
	task.CreatedAt = time.Now()
	tasks[id] = task
	c.JSON(http.StatusOK, nil)
}

func HandleGetAllTasks(c *gin.Context) {
	c.JSON(http.StatusOK, tasks)

}

func HandleToggleComplete(c *gin.Context) {
	idStr := c.Param("taskId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Ты указал неправильный id задачи"+err.Error())
		return
	}

	task := tasks[id]
	task.IsComplete = !task.IsComplete
	tasks[id] = task
	c.JSON(http.StatusOK, task)
}

func HandleGetTask(c *gin.Context) {
	idStr := c.Param("taskId")   // Получаем строку из URL-параметра "taskId", к примеру если запрос пришёл на /tasks/7c2544b2-aa6f..., то idStr = "7c2544b2-aa6f...".
	id, err := uuid.Parse(idStr) // Превращаем строку idStr в тип uuid.UUID
	if err != nil {
		c.JSON(http.StatusBadRequest, "Ты указал неправильный id задачи"+err.Error())
		return
	}

	task, exists := tasks[id] // Находим задачу в мапе tasks по ее id

	if !exists {
		c.JSON(http.StatusNotFound, "Задача не найдена")
		return
	}
	c.JSON(http.StatusOK, task)
}

func HandleDeleteTask(c *gin.Context) {
	idStr := c.Param("taskId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Ты указал неправильный id задачи"+err.Error())
	}

	task, exists := tasks[id]

	if !exists {
		c.JSON(http.StatusNotFound, "Задача не найдена"+err.Error())
	}

	delete(tasks, id) // работает только с мапами

	c.JSON(http.StatusOK, task)
}

func HandleChangeTask(c *gin.Context) {
	idStr := c.Param("taskId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Ты указал неправильный id задачи"+err.Error())
		return
	}

	task, exists := tasks[id]
	if !exists {
		c.JSON(http.StatusNotFound, "Задача не найдена")
		return
	}

	var updates Task
	if err := c.BindJSON(&updates); err != nil { // BindJSON читает тело запроса body и записывает его в структуру updates, тоесть преобразует JSON в Go объект(структуру)
		c.JSON(http.StatusBadRequest, "у тебя неправильный json"+err.Error())
		return
	}

	task.Title = updates.Title
	tasks[id] = task

	c.JSON(http.StatusOK, task)

}
