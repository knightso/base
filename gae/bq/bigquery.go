package bq

import (
	"encoding/json"
	"time"

	"appengine"
	"appengine/taskqueue"

	"github.com/knightso/base/errors"
)

type Task struct {
	LogID    string
	InsertID string
	Time     time.Time
	Record   interface{}
}

func PullReport(c appengine.Context, logID, insertID string, jsonReport interface{}) error {
	v := Task{
		LogID:    logID,
		InsertID: insertID,
		Time:     time.Now(),
		Record:   jsonReport,
	}

	payload, err := json.Marshal(v)
	if err != nil {
		return errors.WrapOr(err)
	}

	task := taskqueue.Task{
		Payload: payload,
		Method:  "PULL",
		Tag:     logID,
	}
	_, err = taskqueue.Add(c, &task, "log2bigquery")
	if err != nil {
		return errors.WrapOr(err)
	}

	return nil
}
