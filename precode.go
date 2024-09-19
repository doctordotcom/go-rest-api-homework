package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

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

// Ниже обработчики для каждого эндпоинта

// Обработчик для получения всех задач
func firstHandle(res http.ResponseWriter, req *http.Request) {
	out, err := json.Marshal(tasks)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(out))
}

// Обработчик для отправки задачи на сервер
func secondHandle(res http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	var task Task
	// прочтение тела запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	// перевод JSON в Task
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	_, ok = tasks[task.ID]
	if ok {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}

// Обработчик для получения задачи по ID
func thirdHandle(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(res, "Задача не найдена.", http.StatusNoContent)
		return
	}

	t, err := json.Marshal(task)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(t)
}

// Обработчик удаления задачи по ID

func deleteTask(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	id := chi.URLParam(req, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(res, "Задача не найдена", http.StatusNoContent)
		return
	}

	delete(tasks, task.ID)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()
	// здесь зарегистрированы обработчики
	r.Get("/tasks", firstHandle)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
