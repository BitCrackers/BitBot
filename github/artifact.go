package github

import "time"

type Artifact struct {
	ID          int       `json:"id"`
	NodeID      string    `json:"node_id"`
	Name        string    `json:"name"`
	SizeInBytes int       `json:"size_in_bytes"`
	Url         string    `json:"url"`
	ArchiveUrl  string    `json:"archive_download_url"`
	Expired     bool      `json:"expired"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}
