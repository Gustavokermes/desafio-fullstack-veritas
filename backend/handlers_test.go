package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type testTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func TestCreateAndGetTask(t *testing.T) {
	// Limpa as tarefas antes do teste
	tasks = nil
	currentID = 1
	saveTasksToFile()

	task := testTask{Title: "Teste", Description: "Desc", Status: "A Fazer"}
	body, _ := json.Marshal(task)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()
	CreateTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("esperado status 201, obteve %d", resp.StatusCode)
	}

	req = httptest.NewRequest("GET", "/tasks", nil)
	w = httptest.NewRecorder()
	GetTasks(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperado status 200, obteve %d", resp.StatusCode)
	}
	var got []Task
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if len(got) != 1 || got[0].Title != "Teste" {
		t.Fatalf("esperado 1 tarefa com título 'Teste', obteve %+v", got)
	}
}

func TestUpdateAndDeleteTask(t *testing.T) {
	tasks = []Task{{ID: "1", Title: "Old", Description: "", Status: Todo}}
	currentID = 2
	saveTasksToFile()

	// Atualizar
	task := testTask{Title: "Novo", Description: "Editado", Status: "Em Progresso"}
	body, _ := json.Marshal(task)
	req := httptest.NewRequest("PUT", "/tasks/1", bytes.NewReader(body))
	vars := map[string]string{"id": "1"}
	req = muxSetVars(req, vars)
	w := httptest.NewRecorder()
	UpdateTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperado status 200, obteve %d", resp.StatusCode)
	}

	// Excluir
	req = httptest.NewRequest("DELETE", "/tasks/1", nil)
	req = muxSetVars(req, vars)
	w = httptest.NewRecorder()
	DeleteTask(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("esperado status 204, obteve %d", resp.StatusCode)
	}
}

func TestCreateTaskInvalid(t *testing.T) {
	tasks = nil
	currentID = 1
	saveTasksToFile()

	// Título vazio
	task := testTask{Title: "", Description: "", Status: "A Fazer"}
	body, _ := json.Marshal(task)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()
	CreateTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("esperado status 400 para título vazio, obteve %d", resp.StatusCode)
	}
}

func TestCreateTaskInvalidStatus(t *testing.T) {
	tasks = nil
	currentID = 1
	saveTasksToFile()

	task := testTask{Title: "Teste", Description: "", Status: "Inexistente"}
	body, _ := json.Marshal(task)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()
	CreateTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("esperado status 201 para status inválido (deve corrigir para 'A Fazer'), obteve %d", resp.StatusCode)
	}
	var got Task
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got.Status != Todo {
		t.Fatalf("esperado status 'A Fazer', obteve %s", got.Status)
	}
}

func TestUpdateTaskNotFound(t *testing.T) {
	tasks = nil
	currentID = 1
	saveTasksToFile()

	task := testTask{Title: "Teste", Description: "", Status: "A Fazer"}
	body, _ := json.Marshal(task)
	req := httptest.NewRequest("PUT", "/tasks/999", bytes.NewReader(body))
	vars := map[string]string{"id": "999"}
	req = muxSetVars(req, vars)
	w := httptest.NewRecorder()
	UpdateTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("esperado status 404 para update de tarefa inexistente, obteve %d", resp.StatusCode)
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	tasks = nil
	currentID = 1
	saveTasksToFile()

	req := httptest.NewRequest("DELETE", "/tasks/999", nil)
	vars := map[string]string{"id": "999"}
	req = muxSetVars(req, vars)
	w := httptest.NewRecorder()
	DeleteTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("esperado status 404 para delete de tarefa inexistente, obteve %d", resp.StatusCode)
	}
}

func TestCreateTaskNoDescription(t *testing.T) {
	tasks = nil
	currentID = 1
	saveTasksToFile()

	task := testTask{Title: "Sem descrição", Status: "A Fazer"}
	body, _ := json.Marshal(task)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()
	CreateTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("esperado status 201, obteve %d", resp.StatusCode)
	}
	var got Task
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got.Description != "" {
		t.Fatalf("esperado descrição vazia, obteve '%s'", got.Description)
	}
}

func TestUpdateTaskInvalidStatus(t *testing.T) {
	tasks = []Task{{ID: "1", Title: "T1", Description: "", Status: Todo}}
	currentID = 2
	saveTasksToFile()

	task := testTask{Title: "T1", Description: "", Status: "Invalido"}
	body, _ := json.Marshal(task)
	req := httptest.NewRequest("PUT", "/tasks/1", bytes.NewReader(body))
	vars := map[string]string{"id": "1"}
	req = muxSetVars(req, vars)
	w := httptest.NewRecorder()
	UpdateTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("esperado status 400 para status inválido, obteve %d", resp.StatusCode)
	}
}

func TestCreateTaskDoneStatus(t *testing.T) {
	tasks = nil
	currentID = 1
	saveTasksToFile()

	task := testTask{Title: "Finalizada", Status: "Concluídas"}
	body, _ := json.Marshal(task)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()
	CreateTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("esperado status 201, obteve %d", resp.StatusCode)
	}
	var got Task
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got.Status != Done {
		t.Fatalf("esperado status 'Concluídas', obteve %s", got.Status)
	}
}

func TestCreateTaskIDGeneration(t *testing.T) {
	tasks = nil
	currentID = 1
	saveTasksToFile()

	task := testTask{Title: "Primeira", Status: "A Fazer"}
	body, _ := json.Marshal(task)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()
	CreateTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("esperado status 201, obteve %d", resp.StatusCode)
	}
	var got Task
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got.ID != "1" {
		t.Fatalf("esperado ID '1', obteve %s", got.ID)
	}
}

// Helper para setar vars do mux
func muxSetVars(r *http.Request, vars map[string]string) *http.Request {
	return mux.SetURLVars(r, vars)
}
