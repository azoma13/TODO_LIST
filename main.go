package main

import (
	"log"
	"net/http"
	"os"

	"github.com/azoma13/go_final_project/configs"

	"github.com/azoma13/go_final_project/internal/dataBase"
	"github.com/azoma13/go_final_project/internal/handlers"

	_ "modernc.org/sqlite"
)

func main() {
	DB := dataBase.DbFunc()
	defer DB.Close()
	port := configs.DefaultPort

	envPort := os.Getenv(configs.ToDoPort)
	if len(envPort) != 0 {
		port = envPort
	}
	port = ":" + port

	serverFiles := http.FileServer(http.Dir(configs.WebDir))
	http.Handle("/", serverFiles)
	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	http.HandleFunc("/api/task", handlers.TaskHandler)
	http.HandleFunc("/api/task/done", handlers.DoneTaskHandler)
	http.HandleFunc("/api/tasks", handlers.GetTasksHandler)

	log.Println("application running on port" + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
