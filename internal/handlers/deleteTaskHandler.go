package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/azoma13/go_final_project/configs"
	"github.com/azoma13/go_final_project/internal/dataBase"
)

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	s, ok := q["id"]
	if !ok {
		sendErrorResponse(w, "missing id url parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(s[0], 10, 32)
	if err != nil {
		sendErrorResponse(w, "failed to convert id to int", http.StatusBadRequest)
		return
	}

	err = dataBase.DeleteTask(int(id))
	if err != nil {
		sendErrorResponse(w, "failed to delete task with id = "+s[0], http.StatusInternalServerError)
		return
	}

	response := configs.CreateTaskResponse{}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
