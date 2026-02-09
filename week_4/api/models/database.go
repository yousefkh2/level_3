package models

import "time"

type CreateDatabaseRequest struct {
	ID        string `json:"id"`
	Instances int    `json:"instances"`
	Storage   string `json:"storage"`
}

type DatabaseSpec struct {
	Instances int    `json:"instances"`
	Storage   string `json:"storage"`
}

type DatabaseResponse struct {
	ID        string       `json:"id"`
	Spec      DatabaseSpec `json:"spec"`
	Status    string       `json:"status"`
	CreatedAt time.Time    `json:"createdAt"`
}
