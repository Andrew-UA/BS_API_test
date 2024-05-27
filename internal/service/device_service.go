package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Andrew-UA/BS_API_test/internal/entity"
	"net/http"
)

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

type TaskRequest struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type AuthService struct {
	Client *http.Client
}

func NewAuthService() *AuthService {
	return &AuthService{
		Client: &http.Client{},
	}
}

func (s *AuthService) Login(host, login, password string) (*entity.Device, error) {
	authReq := AuthRequest{Login: login, Password: password}
	body, err := json.Marshal(authReq)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Post(host+"/api/v0/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed with status code: %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	return entity.NewDevice(host, authResp.AccessToken), nil
}

func (s *AuthService) DoTask(host, token string, task *entity.Task) error {
	taskReq := TaskRequest{Type: task.Type, Payload: task.Payload}
	body, err := json.Marshal(taskReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", host+"/api/v0/do/task", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("task execution failed with status code: %d", resp.StatusCode)
	}

	return nil
}
