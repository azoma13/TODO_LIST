package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/azoma13/go_final_project/configs"
	"github.com/azoma13/go_final_project/internal/dataBase"
)

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var task configs.Task
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		sendErrorResponse(w, "json deserialization error", http.StatusBadRequest)
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

	parseDate, err = time.Parse(configs.DateLayout, task.Date)
	if err != nil {
		sendErrorResponse(w, "date is in an incorrect format", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(task.ID, 10, 32)
	if err != nil {
		sendErrorResponse(w, "error parse int id", http.StatusBadRequest)
		return
	}

	err = dataBase.UpdateTask(int(id), parseDate.Format("20060102"), task.Title, task.Comment, task.Repeat)
	if err != nil {
		sendErrorResponse(w, "failed to update task into db", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(task)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
