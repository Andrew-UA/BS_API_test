package response

import "github.com/Andrew-UA/BS_API_test/internal/dto"

type AllDevicesTaskListResponse struct {
	Items []dto.DeviceDTO `json:"items"`
}
