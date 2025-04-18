package model

import "time"

type Accommodation struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	MainImage         string    `json:"main_image"`
	VillageID         int64     `json:"village_id"`
	About             string    `json:"about"`
	PopularFacilities string    `json:"popular_facilities"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
