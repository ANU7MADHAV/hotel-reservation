package data

import (
	"database/sql"
)

type Models struct {
	Hotels HotelsModel
}

func NewHotelModel(db *sql.DB) Models {
	return Models{
		Hotels: HotelsModel{DB: db},
	}
}
