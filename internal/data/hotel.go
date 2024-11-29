package data

import (
	"database/sql"
	"fmt"
	"time"
)

type HotelsType struct {
	ID        int64     `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Location  string    `json:"location"`
}

type HotelsModel struct {
	DB *sql.DB
}

func (h *HotelsModel) Insert(hotels *HotelsType) error {
	if h.DB == nil {
		return fmt.Errorf("database connection is nil")
	}
	const query = `
	INSERT INTO hotels(name, address, location)
VALUES($1, $2, $3)
RETURNING id;
	`
	args := []interface{}{hotels.Name, hotels.Address, hotels.Location}

	return h.DB.QueryRow(query, args...).Scan(&hotels.ID)
}

func (h *HotelsModel) Get(id int64) (hotels *HotelsType, err error) {
	return nil, nil
}

func (h *HotelsModel) Update(id int64) (hotels *HotelsType, err error) {
	return nil, nil
}

func (h *HotelsModel) Delete(id int64) error {
	return nil
}
