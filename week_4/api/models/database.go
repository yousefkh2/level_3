package models

import "time"

type CreateDatabaseRequest struct {
	Name      string `json:"name"`
	Instances int    `json:"instances"`
	Storage   string `json:"storage"`
}

type DatabaseSpec struct {
	Instances int    `json:"instances"`
	Storage   string `json:"storage"`
}

type DatabaseResponse struct {
	Name      string       `json:"name"`
	Spec      DatabaseSpec `json:"spec"`
	Status    string       `json:"status"`
	CreatedAt time.Time    `json:"createdAt"`
}

type DatabaseDetailResponse struct {
	Name       string         `json:"name"`
	Spec       DatabaseSpec   `json:"spec"`
	Status     string         `json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	Connection ConnectionInfo `json:"connection"`
}

type ConnectionInfo struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateDatabaseRequest struct {
	Instances *int    `json:"instances,omitempty"`
	Storage   *string `json:"storage,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
