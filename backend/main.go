package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	loadTasksFromFile()

	router := mux.NewRouter()

	router.HandleFunc(PathTasks, GetTasks).Methods("GET")
	router.HandleFunc(PathTasks, CreateTask).Methods("POST")
	router.HandleFunc(PathTaskByID, UpdateTask).Methods("PUT")
	router.HandleFunc(PathTaskByID, DeleteTask).Methods("DELETE")
	router.HandleFunc(PathTasks, OptionsHandler).Methods("OPTIONS")
	router.HandleFunc(PathTaskByID, OptionsHandler).Methods("OPTIONS")

	log.Println("🚀 Servidor Kanban rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
