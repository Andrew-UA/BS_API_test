package dto

type DeviceDTO struct {
	ID    string   `json:"device_id"`
	Tasks TasksDTO `json:"tasks"`
}
