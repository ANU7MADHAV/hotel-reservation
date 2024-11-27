package data

import "time"

type Hotels struct {
	ID        int64     `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Location  string    `json:"location"`
}
