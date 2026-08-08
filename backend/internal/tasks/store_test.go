package tasks

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestCreateUpdateDeleteTask(t *testing.T) {
	store, err := NewJSONStore(filepath.Join(t.TempDir(), "tasks.json"))
	if err != nil {
		t.Fatalf("NewJSONStore() error = %v", err)
	}

	created, err := store.Create(TaskPayload{
		Title:       "Preparar README",
		Description: "Documentar como rodar o projeto",
		Status:      StatusTodo,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.ID == "" {
		t.Fatal("Create() returned empty id")
	}

	updated, err := store.Update(created.ID, TaskPayload{
		Title:       "Preparar README",
		Description: "Documentar endpoints e decisoes tecnicas",
		Status:      StatusInProgress,
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.Status != StatusInProgress {
		t.Fatalf("Update() status = %q, want %q", updated.Status, StatusInProgress)
	}

	if err := store.Delete(created.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if len(store.List()) != 0 {
		t.Fatal("List() should be empty after delete")
	}
}

func TestCreateValidatesTitle(t *testing.T) {
	store, err := NewJSONStore(filepath.Join(t.TempDir(), "tasks.json"))
	if err != nil {
		t.Fatalf("NewJSONStore() error = %v", err)
	}

	_, err = store.Create(TaskPayload{Title: "   ", Status: StatusTodo})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("Create() error = %v, want ErrValidation", err)
	}
}
