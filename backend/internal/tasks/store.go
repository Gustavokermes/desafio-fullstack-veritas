package tasks

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ErrNotFound   = errors.New("task not found")
	ErrValidation = errors.New("invalid task")
)

type JSONStore struct {
	mu    sync.RWMutex
	path  string
	tasks map[string]Task
}

func NewJSONStore(path string) (*JSONStore, error) {
	store := &JSONStore{
		path:  path,
		tasks: make(map[string]Task),
	}

	if err := store.load(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *JSONStore) List() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
	})

	return tasks
}

func (s *JSONStore) Create(payload TaskPayload) (Task, error) {
	title, description, status, err := normalizePayload(payload, StatusTodo)
	if err != nil {
		return Task{}, err
	}

	now := time.Now().UTC()
	task := Task{
		ID:          newID(),
		Title:       title,
		Description: description,
		Status:      status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks[task.ID] = task
	if err := s.saveLocked(); err != nil {
		delete(s.tasks, task.ID)
		return Task{}, err
	}

	return task, nil
}

func (s *JSONStore) Update(id string, payload TaskPayload) (Task, error) {
	title, description, status, err := normalizePayload(payload, StatusTodo)
	if err != nil {
		return Task{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}

	task.Title = title
	task.Description = description
	task.Status = status
	task.UpdatedAt = time.Now().UTC()

	s.tasks[id] = task
	if err := s.saveLocked(); err != nil {
		return Task{}, err
	}

	return task, nil
}

func (s *JSONStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrNotFound
	}

	delete(s.tasks, id)
	return s.saveLocked()
}

func (s *JSONStore) load() error {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return os.MkdirAll(filepath.Dir(s.path), 0o755)
	}
	if err != nil {
		return err
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return fmt.Errorf("invalid tasks JSON: %w", err)
	}

	for _, task := range tasks {
		if task.ID == "" {
			continue
		}
		s.tasks[task.ID] = task
	}

	return nil
}

func (s *JSONStore) saveLocked() error {
	tasks := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
	})

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0o644)
}

func normalizePayload(payload TaskPayload, defaultStatus Status) (string, string, Status, error) {
	title := strings.TrimSpace(payload.Title)
	description := strings.TrimSpace(payload.Description)
	status := payload.Status
	if status == "" {
		status = defaultStatus
	}

	if title == "" {
		return "", "", "", fmt.Errorf("%w: title is required", ErrValidation)
	}
	if len(title) > 80 {
		return "", "", "", fmt.Errorf("%w: title must have at most 80 characters", ErrValidation)
	}
	if len(description) > 280 {
		return "", "", "", fmt.Errorf("%w: description must have at most 280 characters", ErrValidation)
	}
	if !validStatus(status) {
		return "", "", "", fmt.Errorf("%w: unsupported status", ErrValidation)
	}

	return title, description, status, nil
}

func newID() string {
	buffer := make([]byte, 8)
	if _, err := rand.Read(buffer); err == nil {
		return hex.EncodeToString(buffer)
	}
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
