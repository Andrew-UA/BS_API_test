package service

import (
	"errors"
	"fmt"
	"github.com/Andrew-UA/BS_API_test/internal/dto"
	"github.com/Andrew-UA/BS_API_test/internal/entity"
	"log/slog"
	"time"
)

var ErrNotFound = errors.New("not found device")

type IDeviceService interface {
	Login(host string, login string, password string) (*entity.Device, error)
	DoTask(host string, token string, task *entity.Task) error
}

type TaskManagerService struct {
	deviceService IDeviceService
	Devices       map[string]*entity.Device
}

func NewTaskManagerService(deviceService IDeviceService) *TaskManagerService {
	return &TaskManagerService{
		deviceService: deviceService,
		Devices:       make(map[string]*entity.Device),
	}
}

func (tm *TaskManagerService) AddTask(deviceId string, task entity.Task) error {
	device, ok := tm.Devices[deviceId]
	if !ok {
		return ErrNotFound
	}

	return device.AddTask(&task)
}

func (tm *TaskManagerService) AddDevice(host string, login string, password string) error {
	device, err := tm.deviceService.Login(host, login, password)
	if err != nil {
		return err
	}

	tm.Devices[device.ID] = device

	go tm.processDevice(device)

	return nil
}

func (tm *TaskManagerService) GetAllTaskList() *[]dto.DeviceDTO {
	dtoList := make([]dto.DeviceDTO, len(tm.Devices), len(tm.Devices))

	i := 0
	for _, device := range tm.Devices {
		dtoList[i] = dto.DeviceDTO{
			ID: device.ID,
			Tasks: dto.TasksDTO{
				Queued:    device.QueuedTasksCount(),
				Processed: device.ProcessedTasksCount(),
			},
		}
		i++
	}

	return &dtoList
}

func (tm *TaskManagerService) GetDeviceTaskList(deviceId string) (*dto.TasksDTO, error) {
	device, ok := tm.Devices[deviceId]
	if !ok {
		return nil, fmt.Errorf("device not found")
	}

	return &dto.TasksDTO{
		Queued:    device.QueuedTasksCount(),
		Processed: device.ProcessedTasksCount(),
	}, nil
}

func (tm *TaskManagerService) ClearTasks(deviceId string) error {
	device, ok := tm.Devices[deviceId]
	if !ok {
		return ErrNotFound
	}

	previousStatus := device.Status

	device.ClearQueued()

	if previousStatus == entity.STOPPED {
		go tm.processDevice(device)
	}

	return nil
}

func (tm *TaskManagerService) Work() {
	for _, device := range tm.Devices {
		go tm.processDevice(device)
	}
}

func (tm *TaskManagerService) processDevice(device *entity.Device) {
	for {
		task := device.ShiftTask()
		if task == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		err := tm.deviceService.DoTask(device.Host, device.AccessToken, task)
		if err != nil {
			device.StopProcessing()
			slog.Error("Processing task error", "err", err.Error(), "device_id", device.ID, "host", device.Host)
			return
		}

		device.IncrementProcessedTasks()
	}
}
