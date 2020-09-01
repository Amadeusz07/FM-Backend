package dtos

import "time"

type LoginResponse struct {
	Token       string    `json:"token"`
	ExpiresDate time.Time `json:"expiresDate"`
	Username    string    `json:"username"`
	ProjectId   string    `json:"projectId,omitempty"`
}
