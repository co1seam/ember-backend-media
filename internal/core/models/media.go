package models

import "time"

type Media struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ContentType string    `json:"content_type"`
	StoragePath string    `json:"storage_path"`
	OwnerID     string    `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	URL         string    `json:"url"`
}

type CreateMediaRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ContentType string `json:"content_type"`
	OwnerID     string `json:"owner_id"`
}

type UpdateMediaRequest struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GetMediaRequest struct {
	ID string `json:"id"`
}

type ListMediaRequest struct {
	OwnerID string `json:"owner_id"`
}
