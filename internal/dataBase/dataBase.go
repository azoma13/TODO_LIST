package dataBase

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/azoma13/go_final_project/configs"
)

var DB *sql.DB

func DbFunc() *sql.DB {

	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")

	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	DB, err = sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Fatal(err)
	}

	if install {
		_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat VARCHAR(128) NOT NULL
			);`)
		if err != nil {
			log.Fatal(err)
		}
		_, err = DB.Exec(`CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);`)
		if err != nil {
			log.Fatal(err)
		}
	}
	return DB
}

func AddTask(date, title, comment, repeat string) (string, error) {

	result, err := DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", date, title, comment, repeat)
	if err != nil {
		return "", fmt.Errorf("error execute in func AddTask: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("error while getting lasi inserted id: %w", err)
	}

	return strconv.FormatInt(id, 10), nil
}

func GetAllTasks() ([]configs.Task, error) {

	rows, err := DB.Query("SELECT * FROM scheduler ORDER BY date ASC LIMIT 50")
	if err != nil {
		return nil, fmt.Errorf("error query in func GetAllTasks: %w", err)
	}
	defer rows.Close()

	var tasks []configs.Task
	for rows.Next() {
		task := configs.Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("error scan in func GetAllTasks: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error cursor rows in func GetAllTasks: %w", err)
	}

	if tasks == nil {
		tasks = []configs.Task{}
	}

	return tasks, nil
}

func SearchTaskById(id int64) (*configs.Task, error) {

	row := DB.QueryRow("SELECT * FROM scheduler where id = ?", id)

	var task configs.Task

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, fmt.Errorf("error scan in func SearchTaskById: %w", err)
	}

	return &task, nil
}

func UpdateTask(id int, date string, title, comment, repeat string) error {

	sqlQuery := `
	UPDATE scheduler
	SET date = ?, title = ?, comment = ?, repeat = ? 
	WHERE id = ?`
	_, err := DB.Exec(sqlQuery, date, title, comment, repeat, id)
	if err != nil {
		return fmt.Errorf("error execute in func UpdateTask: %w", err)
	}

	return nil
}

func DeleteTask(id int) error {

	_, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error execute in func DeleteTask: %w", err)
	}

	return nil
}
