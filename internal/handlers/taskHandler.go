package handlers

import "net/http"

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		AddTaskHandler(w, r)
	} else if r.Method == "GET" {
		GetTaskByIdHandler(w, r)
	} else if r.Method == "PUT" {
		UpdateTaskHandler(w, r)
	} else if r.Method == "DELETE" {
		DeleteTaskHandler(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
