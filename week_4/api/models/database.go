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
	Connection ConnectionInfo `json:"connection"` // NEW
}

type ConnectionInfo struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
}
