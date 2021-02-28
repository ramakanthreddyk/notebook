package models

type ToDoList struct {
	TaskId *int               `json:"task_id"`
	Task   string             `json:"task,omitempty"`
}
