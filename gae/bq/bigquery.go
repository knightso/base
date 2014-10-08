package bq

import (
	"time"
)

type Task struct {
	LogID    string
	InsertID string
	Time     time.Time
	Record   interface{}
}
