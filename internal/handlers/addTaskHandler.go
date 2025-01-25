package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/azoma13/go_final_project/configs"
	"github.com/azoma13/go_final_project/internal/dataBase"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var task configs.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		sendErrorResponse(w, "json deserialization Error", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		sendErrorResponse(w, "task title is not correct", http.StatusBadRequest)
		return
	}

	now := time.Now()
	now = now.Truncate(24 * time.Hour)
	today := now.Format(configs.DateLayout)

	if task.Date == "" {
		task.Date = today
	}

	parseDate, err := time.Parse(configs.DateLayout, task.Date)
	if err != nil {
		sendErrorResponse(w, "date is in an incorrect format", http.StatusBadRequest)
		return
	}

	if parseDate.Before(now) {
		if task.Repeat == "" {
			task.Date = today
		} else {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				sendErrorResponse(w, "error in repetition rule", http.StatusInternalServerError)
				return
			}

			if task.Repeat == "d 1" && nextDate == today {
				task.Date = today
			} else {
				task.Date = nextDate
			}
		}
	}

	if task.Repeat != "" {
		_, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			sendErrorResponse(w, "invalid repetition rule", http.StatusInternalServerError)
			return
		}
	}

	parseDate, err = time.Parse(configs.DateLayout, task.Date)
	if err != nil {
		sendErrorResponse(w, "date is in an incorrect format", http.StatusBadRequest)
		return
	}

	id, err := dataBase.AddTask(parseDate.Format("20060102"), task.Title, task.Comment, task.Repeat)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

func sendErrorResponse(w http.ResponseWriter, errorMessage string, statusServer int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusServer)
	json.NewEncoder(w).Encode(map[string]string{"error": errorMessage})
}
