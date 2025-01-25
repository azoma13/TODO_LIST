package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/azoma13/go_final_project/internal/dataBase"
)

func GetTaskByIdHandler(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	s, ok := q["id"]
	if !ok {
		sendErrorResponse(w, "missing id url parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(s[0], 10, 64)
	if err != nil {
		sendErrorResponse(w, "failed to convert id to int", http.StatusBadRequest)
		return
	}

	task, err := dataBase.SearchTaskById(id)
	if err != nil {
		sendErrorResponse(w, "task with id = "+strconv.Itoa(int(id))+" not found", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(task)
	if err != nil {
		sendErrorResponse(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}
