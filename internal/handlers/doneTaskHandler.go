package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/azoma13/go_final_project/configs"
	"github.com/azoma13/go_final_project/internal/dataBase"
)

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {

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

	task, err := dataBase.SearchTaskById(id)
	if err != nil {
		sendErrorResponse(w, "failed to search task by id", http.StatusInternalServerError)
		return
	}

	if task.Repeat == "" {
		err = dataBase.DeleteTask(int(id))
		if err != nil {
			sendErrorResponse(w, "failed to delete task with id = "+s[0], http.StatusInternalServerError)
			return
		}

	} else {

		date, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			sendErrorResponse(w, "failed to get next date", http.StatusInternalServerError)
			return
		}

		err = dataBase.UpdateTask(int(id), date, task.Title, task.Comment, task.Repeat)
		if err != nil {
			sendErrorResponse(w, "failed to update task into db", http.StatusInternalServerError)
			return
		}
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
