package v1

import "time"

type GrowthRecord struct {
	Date time.Time `json:"date"`
	Label string `json:"label"`
}
