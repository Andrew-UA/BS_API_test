package http

import (
	"encoding/json"
	"errors"
	"github.com/Andrew-UA/BS_API_test/internal/api/http/request"
	"github.com/Andrew-UA/BS_API_test/internal/api/http/response"
	"github.com/Andrew-UA/BS_API_test/internal/service"
	"net/http"
)

type TaskHandler struct {
	taskManagerService *service.TaskManagerService
}

func NewTaskHandler(t *service.TaskManagerService) *TaskHandler {
	return &TaskHandler{
		taskManagerService: t,
	}
}

func (t *TaskHandler) HandelQueueTask(w http.ResponseWriter, r *http.Request) error {
	var createTaskRequest request.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&createTaskRequest); err != nil {
		return APIError{
			StatusCode: http.StatusBadRequest,
			Message:    "Bad Request",
		}
	}

	if err := t.taskManagerService.AddTask(createTaskRequest.DeviceID, createTaskRequest.Task); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return NewAPIError(http.StatusNotFound, err)
		}
		return err
	}

	writeJSON(w, http.StatusCreated, make(map[string]interface{}))

	return nil
}
func (t *TaskHandler) HandleAllDevicesTaskList(w http.ResponseWriter, r *http.Request) error {
	writeJSON(w, 200,
		response.AllDevicesTaskListResponse{
			Items: *t.taskManagerService.GetAllTaskList(),
		},
	)

	return nil
}

func (t *TaskHandler) HandelDeviceTaskList(w http.ResponseWriter, r *http.Request) error {
	deviceID := r.PathValue("device_id")
	taskDTO, err := t.taskManagerService.GetDeviceTaskList(deviceID)
	if err != nil {
		return APIError{
			StatusCode: http.StatusNotFound,
			Message:    "Device not found",
		}
	}
	writeJSON(w, 200, taskDTO)

	return nil
}

func (t *TaskHandler) HandleDeviceTaskClear(w http.ResponseWriter, r *http.Request) error {
	deviceID := r.PathValue("device_id")
	if err := t.taskManagerService.ClearTasks(deviceID); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return NewAPIError(http.StatusNotFound, err)
		}
		return err
	}

	taskDTO, err := t.taskManagerService.GetDeviceTaskList(deviceID)
	if err != nil {
		return APIError{
			StatusCode: http.StatusNotFound,
			Message:    "Device not found",
		}
	}
	writeJSON(w, 200, taskDTO)

	return nil
}
