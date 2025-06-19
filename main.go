package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Task struct {
	Id          uuid.UUID `json: "id"`
	Title       string    `json: "title"`
	Description string    `json: "description"`
	IsComplete  bool      `json: "is_complete"`
	CreatedAt   time.Time `json: "created_at"`
}

var tasks = make(map[uuid.UUID]Task)

func main() {

	r := gin.Default() // роутер который направляет запросы куда надо, для него мы создаем пути в коде снизу

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

	// r.POST("/tasks", func(c *gin.Context) { // эти пути называются эндпоинты
	// 	c.JSON(http.StatusOK, gin.H{ // тут мы возвращаем клиенту JSON запрос который он сделал, пример {"message": "json"}
	// 		"message": "hello get",
	// 	})
	// })

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
