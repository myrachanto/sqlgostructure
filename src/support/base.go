 package support

import (
	"time"
)

type Base struct {
	Created_At time.Time  `bson:"created_at"`
	Updated_At time.Time  `bson:"updated_at"`
	Delete_At  *time.Time `bson:"deleted_at"`
}

type Location struct {
	Longtitude float64 `json:"longtitude,omitempty"`
	Latitude   float64 `json:"latitude,omitempty"`
}
