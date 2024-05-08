package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5" // Импорт пакета chi для маршрутизации
)

// Структура для задачи
type Task struct {
	ID           string   `json:"id"`           // ID задачи
	Description  string   `json:"description"`  // Описание задачи
	Note         string   `json:"note"`         // Примечание к задаче
	Applications []string `json:"applications"` // Приложения, используемые для выполнения задачи
}

// Мапа для хранения задач
var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func generateNewID() string {
	return strconv.Itoa(len(tasks) + 1)
}

// Обработчик для получения всех задач
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Преобразование мапы задач в массив
	var tasksArray []Task
	for _, task := range tasks {
		tasksArray = append(tasksArray, task)
	}

	err := json.NewEncoder(w).Encode(tasksArray)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Обработчик для создания новой задачи
func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Если ID задачи не указан, генерируем новый
	if newTask.ID == "" {
		newTask.ID = generateNewID()
	}

	// Если Applications не указаны, добавляем User-Agent из запроса
	if len(newTask.Applications) == 0 {
		newTask.Applications = []string{r.UserAgent()}
	}

	// Проверяем, существует ли уже задача с таким ID
	_, ok := tasks[newTask.ID]
	if ok {
		http.Error(w, "Задача с таким ID уже существует", http.StatusBadRequest)
		return
	}

	tasks[newTask.ID] = newTask
	w.WriteHeader(http.StatusCreated)
}

// Обработчик для получения задачи по ID
func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Некорректный запрос серверу", http.StatusBadRequest)
		return
	}

	err := json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Обработчик для удаления задачи по ID
func deleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}
	delete(tasks, id)
	w.WriteHeader(http.StatusOK)
}

// Главная функция, запускающая сервер и регистрирующая обработчики
func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getTasks)           // Регистрация обработчика для получения всех задач
	r.Post("/tasks", createTask)        // Регистрация обработчика для создания новой задачи
	r.Get("/tasks/{id}", getTask)       // Регистрация обработчика для получения задачи по ID
	r.Delete("/tasks/{id}", deleteTask) // Регистрация обработчика для удаления задачи по ID

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
