package entity

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
	"sync/atomic"
)

type ProcessingStatus string

const PROCESSING = ProcessingStatus("processing")
const STOPPED = ProcessingStatus("stopped")

type Device struct {
	ID             string
	Host           string
	AccessToken    string
	Status         ProcessingStatus
	QueuedTasks    []*Task
	ProcessedTasks int32
	mx             sync.Mutex
}

func NewDevice(host string, accessToken string) *Device {
	return &Device{
		ID:             uuid.New().String()[:8],
		Host:           host,
		AccessToken:    accessToken,
		Status:         PROCESSING,
		QueuedTasks:    make([]*Task, 0, 100),
		ProcessedTasks: 0,
	}
}

func (d *Device) QueuedTasksCount() int {
	return len(d.QueuedTasks)
}

func (d *Device) ProcessedTasksCount() int {
	return int(d.ProcessedTasks)
}

func (d *Device) AddTask(task *Task) error {
	d.mx.Lock()
	defer d.mx.Unlock()

	if d.QueuedTasksCount() >= 5 {
		return fmt.Errorf("task buffer limit reached for device: %s ", d.ID)
	}
	if d.Status == STOPPED {
		return fmt.Errorf("task processing stopped for device: %s ", d.ID)
	}

	d.QueuedTasks = append(d.QueuedTasks, task)

	return nil
}

func (d *Device) ShiftTask() *Task {
	d.mx.Lock()
	defer d.mx.Unlock()

	if d.QueuedTasksCount() > 0 {
		task := d.QueuedTasks[0]
		d.QueuedTasks = d.QueuedTasks[1:]

		return task
	}

	return nil
}

func (d *Device) StopProcessing() {
	d.mx.Lock()
	defer d.mx.Unlock()

	d.Status = STOPPED
}

func (d *Device) ClearQueued() {
	d.mx.Lock()
	defer d.mx.Unlock()

	d.QueuedTasks = make([]*Task, 0, 100)
	d.Status = PROCESSING
}

func (d *Device) IncrementProcessedTasks() {
	atomic.AddInt32(&d.ProcessedTasks, 1)
}
