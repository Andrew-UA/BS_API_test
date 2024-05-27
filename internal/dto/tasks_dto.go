package dto

type TasksDTO struct {
	Queued    int `json:"queued"`
	Processed int `json:"processed"`
}
