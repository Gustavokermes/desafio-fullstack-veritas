package tasks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestHandlerCreateUpdateMoveAndDelete(t *testing.T) {
	handler, store := newTestHandler(t)

	created := sendTaskRequest(t, handler, http.MethodPost, "/tasks", TaskPayload{
		Title:       "Teste A",
		Description: "Primeira tarefa",
		Status:      StatusTodo,
	}, http.StatusCreated)

	if len(store.List()) != 1 {
		t.Fatalf("Create stored %d tasks, want 1", len(store.List()))
	}

	updated := sendTaskRequest(t, handler, http.MethodPut, "/tasks/"+created.ID, TaskPayload{
		Title:       "Teste A editado",
		Description: "Primeira tarefa editada",
		Status:      StatusTodo,
	}, http.StatusOK)

	if updated.ID != created.ID {
		t.Fatalf("Update changed id from %q to %q", created.ID, updated.ID)
	}
	if updated.Title != "Teste A editado" {
		t.Fatalf("Update title = %q, want %q", updated.Title, "Teste A editado")
	}
	if len(store.List()) != 1 {
		t.Fatalf("Update stored %d tasks, want 1", len(store.List()))
	}

	moved := sendTaskRequest(t, handler, http.MethodPut, "/tasks/"+created.ID, TaskPayload{
		Title:       updated.Title,
		Description: updated.Description,
		Status:      StatusInProgress,
	}, http.StatusOK)

	if moved.Status != StatusInProgress {
		t.Fatalf("Move status = %q, want %q", moved.Status, StatusInProgress)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "/tasks/"+created.ID, nil)
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("DELETE status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
	if len(store.List()) != 0 {
		t.Fatalf("Delete stored %d tasks, want 0", len(store.List()))
	}
}

func TestHandlerCreatesStableUniqueIDs(t *testing.T) {
	handler, _ := newTestHandler(t)

	first := sendTaskRequest(t, handler, http.MethodPost, "/tasks", TaskPayload{Title: "Mesmo titulo"}, http.StatusCreated)
	second := sendTaskRequest(t, handler, http.MethodPost, "/tasks", TaskPayload{Title: "Mesmo titulo"}, http.StatusCreated)

	if first.ID == "" || second.ID == "" {
		t.Fatal("Create returned an empty id")
	}
	if first.ID == second.ID {
		t.Fatalf("Create reused id %q", first.ID)
	}
}

func TestHandlerRejectsInvalidPayloads(t *testing.T) {
	handler, _ := newTestHandler(t)

	tests := []struct {
		name string
		body string
	}{
		{
			name: "missing title",
			body: `{"description":"sem titulo","status":"todo"}`,
		},
		{
			name: "blank title",
			body: `{"title":"   ","status":"todo"}`,
		},
		{
			name: "invalid status",
			body: `{"title":"Teste","status":"blocked"}`,
		},
		{
			name: "invalid json",
			body: `{"title":`,
		},
		{
			name: "trailing json",
			body: `{"title":"Teste"}{"title":"Outro"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(tt.body))
			handler.ServeHTTP(recorder, request)

			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("POST status = %d, want %d", recorder.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestHandlerReturnsNotFoundForMissingTask(t *testing.T) {
	handler, _ := newTestHandler(t)

	sendRawRequest(t, handler, http.MethodPut, "/tasks/does-not-exist", `{"title":"Teste","status":"todo"}`, http.StatusNotFound)
	sendRawRequest(t, handler, http.MethodDelete, "/tasks/does-not-exist", "", http.StatusNotFound)
}

func TestHandlerRejectsUnsupportedMethods(t *testing.T) {
	handler, _ := newTestHandler(t)

	sendRawRequest(t, handler, http.MethodPatch, "/tasks", "", http.StatusMethodNotAllowed)
	sendRawRequest(t, handler, http.MethodPatch, "/tasks/seed-1", "", http.StatusMethodNotAllowed)
}

func newTestHandler(t *testing.T) (*Handler, *JSONStore) {
	t.Helper()

	store, err := NewJSONStore(filepath.Join(t.TempDir(), "tasks.json"))
	if err != nil {
		t.Fatalf("NewJSONStore() error = %v", err)
	}

	return NewHandler(store), store
}

func sendTaskRequest(t *testing.T, handler http.Handler, method, path string, payload TaskPayload, wantStatus int) Task {
	t.Helper()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, path, bytes.NewReader(body))
	handler.ServeHTTP(recorder, request)

	if recorder.Code != wantStatus {
		t.Fatalf("%s %s status = %d, want %d. Body: %s", method, path, recorder.Code, wantStatus, recorder.Body.String())
	}

	var task Task
	if err := json.NewDecoder(recorder.Body).Decode(&task); err != nil {
		t.Fatalf("Decode response error = %v. Body: %s", err, recorder.Body.String())
	}

	return task
}

func sendRawRequest(t *testing.T, handler http.Handler, method, path, body string, wantStatus int) {
	t.Helper()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	handler.ServeHTTP(recorder, request)

	if recorder.Code != wantStatus {
		t.Fatalf("%s %s status = %d, want %d. Body: %s", method, path, recorder.Code, wantStatus, recorder.Body.String())
	}
}
