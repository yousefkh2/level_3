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
