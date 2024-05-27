package http

import "net/http"

func NewRouter(taskHandler *TaskHandler) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /api/v0/devices/tasks/queue", MakeHandler(taskHandler.HandelQueueTask))
	router.HandleFunc("GET /api/v0/devices/tasks/list", MakeHandler(taskHandler.HandleAllDevicesTaskList))
	router.HandleFunc("GET /api/v0/devices/{device_id}/tasks/list", MakeHandler(taskHandler.HandelDeviceTaskList))
	router.HandleFunc("POST /api/v0/devices/{device_id}/tasks/clear", MakeHandler(taskHandler.HandleDeviceTaskClear))

	return router
}
