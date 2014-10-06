package bq

import (
	"time"
)

type Task struct {
	LogID    string
	InsertID string
	Date     time.Time
	Record   interface{}
}
