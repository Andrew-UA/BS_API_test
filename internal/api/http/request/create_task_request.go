package request

import "github.com/Andrew-UA/BS_API_test/internal/entity"

type CreateTaskRequest struct {
	DeviceID string `json:"device_id"`
	Task     entity.Task
}
