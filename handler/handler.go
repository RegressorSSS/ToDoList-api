package handler

import (
	"net/http"

	"github.com/RegressorSSS/model"
	"github.com/RegressorSSS/todolist/model"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Handler struct { // Хранит подключение к БД, и даёт нам возможность создавать HTTP - запросы, которые могут работать с базой
	db *pgx.Conn // Указатель на соединение с БД, через которое мы выполняем SQL - запросы
}

func New(db *pgx.Conn) *Handler { // Конструктор New создаёт новый экземпляр Handler и передаёт ему готовое подключение к БД.
	return &Handler{ // Создаём новый объект Handler и кладём это соединение внутрь Handler.db
		db: db, // первое db — это имя поля структуры, а второе — переменная, содержащая уже установленное подключение к БД
	}
}

func (h *Handler) HandlerCreateTask(c *gin.Context) { // h это экземплер структуры Handler, и он имеет доступ к БД через h.db
	var task model.Task
	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У тебя неправильный json"+err.Error())
		return
	}

	_, err = h.db.Exec(c, "INSERT INTO tasks (title, description) VALUES ($1,$2)", task.Title, task.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "ошибка при сохранении задачи"+err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *Handler) HandleGetAllTasks(c *gin.Context) {
	rows, err := h.db.Query(c, "SELECT * FROM tasks") // Query тут выполняет SQL запрос: получить все строки из таблицы tasks
	if err != nil {
		c.JSON(http.StatusInternalServerError, "ошибка при получении задачи из БД"+err.Error())
		return
	}

	var tasks []model.Task

	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.IsComplete, &task.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "ошибка при Scan из БД"+err.Error())
			return
		}

		tasks = append(tasks, task)
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) HandleToggleComplete(c *gin.Context) {
	// TODO: реализовать
}

func (h *Handler) HandleGetTaskByID(c *gin.Context) {
	// TODO: реализовать
}

func (h *Handler) HandleDeleteTask(c *gin.Context) {
	// TODO: реализовать
}

func (h *Handler) HandleUpdateTask(c *gin.Context) {
	// TODO: реализовать
}
