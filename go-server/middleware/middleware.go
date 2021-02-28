package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-server/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var mainDB *sql.DB

// create connection with sqlite db
func init() {
	createDBInstance()
}

func createDBInstance() {
	db, _ := sql.Open("sqlite3", "notes.db")

	mainDB = db
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS notes(task_id INTEGER PRIMARY KEY AUTOINCREMENT, task varchar(250))")
	if err != nil {
		log.Fatal(err)
	}
	_, sqlerr := stmt.Exec()
	if sqlerr != nil {
		log.Fatal(sqlerr)
	}
}

// GetAllTask get all the task route
func GetAllTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rows, err := mainDB.Query("SELECT * FROM notes")
	var notes []models.ToDoList
	for rows.Next() {
		var note models.ToDoList
		err = rows.Scan(&note.TaskId, &note.Task)
		if err != nil {
			log.Fatal(err)
		}
		notes = append(notes, note)
	}
	json.NewEncoder(w).Encode(notes)
}

// CreateTask create task route
func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var task models.ToDoList
	_ = json.NewDecoder(r.Body).Decode(&task)
	fmt.Println(task)
	insertOneTask(task)
	json.NewEncoder(w).Encode(task)
}

// DeleteTask delete one task route
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	params := mux.Vars(r)
	deleteOneTask(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func deleteOneTask(id string) {
	stmt, err := mainDB.Prepare("DELETE FROM notes WHERE task_id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, er := stmt.Exec(id)

	if er != nil {
		log.Fatal(er)
	}
}

// Insert one task in the DB
func insertOneTask(task models.ToDoList) {
	stmt, err := mainDB.Prepare("INSERT INTO notes('task') values (?)")

	if err != nil {
		log.Fatal(err)
	}

	_, er := stmt.Exec(task.Task)

	if er != nil {
		log.Fatal(er)
	}
}
