package main

// Constantes para evitar duplicação de literais
const (
	PathTasks         = "/tasks"
	PathTaskByID      = "/tasks/{id}"
	HeaderContentType = "Content-Type"
	ContentTypeJSON   = "application/json"
)

type Status string

const (
	Todo       Status = "A Fazer"
	InProgress Status = "Em Progresso"
	Done       Status = "Concluídas"
)

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}

var validStatuses = map[Status]bool{
	Todo:       true,
	InProgress: true,
	Done:       true,
}
