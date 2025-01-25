package configs

const (
	DateLayout  = "20060102"
	DefaultPort = "7540"
	WebDir      = "./web"
	ToDoPort    = "TODO_PORT"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
	Error string `json:"error,omitempty"`
}

type CreateTaskResponse struct {
	Id    int    `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}
