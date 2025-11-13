package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

const tasksFile = "tasks.json"

var (
	tasks     []Task
	currentID int = 1
	mu        sync.RWMutex
)

func loadTasksFromFile() {
	file, err := os.Open(tasksFile)
	if err != nil {
		return // Se não existir, começa vazio
	}
	defer file.Close()
	var loaded []Task
	if err := json.NewDecoder(file).Decode(&loaded); err == nil {
		tasks = loaded
		// Atualiza o currentID para o maior ID + 1
		maxID := 0
		for _, t := range tasks {
			if id, err := strconv.Atoi(t.ID); err == nil && id > maxID {
				maxID = id
			}
		}
		currentID = maxID + 1
	}
}

func saveTasksToFile() {
	file, err := os.Create(tasksFile)
	if err != nil {
		return
	}
	defer file.Close()
	_ = json.NewEncoder(file).Encode(tasks)
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", HeaderContentType)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	w.Header().Set(HeaderContentType, ContentTypeJSON)
	json.NewEncoder(w).Encode(tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	w.Header().Set(HeaderContentType, ContentTypeJSON)

	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		http.Error(w, "Título é obrigatório", http.StatusBadRequest)
		return
	}

	if !validStatuses[task.Status] {
		task.Status = Todo
	}

	mu.Lock()
	task.ID = strconv.Itoa(currentID)
	currentID++
	tasks = append(tasks, task)
	saveTasksToFile()
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	w.Header().Set(HeaderContentType, ContentTypeJSON)

	params := mux.Vars(r)
	taskID := params["id"]

	var updatedTask Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(updatedTask.Title) == "" {
		http.Error(w, "Título é obrigatório", http.StatusBadRequest)
		return
	}

	if !validStatuses[updatedTask.Status] {
		http.Error(w, "Status inválido", http.StatusBadRequest)
		return
	}

	mu.Lock()
	for i, task := range tasks {
		if task.ID == taskID {
			updatedTask.ID = taskID
			tasks[i] = updatedTask
			saveTasksToFile()
			mu.Unlock()
			json.NewEncoder(w).Encode(updatedTask)
			return
		}
	}
	mu.Unlock()

	http.Error(w, "Task não encontrada", http.StatusNotFound)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	params := mux.Vars(r)
	taskID := params["id"]

	mu.Lock()
	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			saveTasksToFile()
			mu.Unlock()
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	mu.Unlock()

	http.Error(w, "Task não encontrada", http.StatusNotFound)
}

func OptionsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	w.WriteHeader(http.StatusOK)
}
