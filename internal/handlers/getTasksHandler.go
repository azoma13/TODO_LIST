package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/azoma13/go_final_project/configs"
	"github.com/azoma13/go_final_project/internal/dataBase"
)

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {

	tasks, err := dataBase.GetAllTasks()
	if err != nil {
		sendErrorResponse(w, "failed to load tasks from db", http.StatusInternalServerError)
		return
	}

	response := configs.TasksResponse{}
	response.Tasks = tasks

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		sendErrorResponse(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
