package service

import (
	"fmt"
	"github.com/Andrew-UA/BS_API_test/internal/entity"
	"math/rand"
	"time"
)

type DeviceServiceMock struct {
}

func (d *DeviceServiceMock) Login(host string, login string, password string) (*entity.Device, error) {

	rToken := randStringBytes(10)

	return entity.NewDevice(host, rToken), nil
}

func (d *DeviceServiceMock) DoTask(host string, token string, task *entity.Task) error {

	randomNumber := rand.Intn(10) + 1

	if randomNumber == 1 {
		return fmt.Errorf("some device error")
	}

	time.Sleep(5 * time.Second)
	return nil

}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
