package bq

import (
	"encoding/json"
	"time"

	"appengine"
	"appengine/taskqueue"

	"github.com/knightso/base/errors"
)

var EnableLog bool = false // debugging variable.

const (
	QUEUE_NAME string = "log2bigquery"
)

type Task struct {
	LogID    string
	InsertID string
	Time     time.Time
	Record   interface{}
}

func SendLog(c appengine.Context, logID, insertID string, record interface{}) error {
	if EnableLog == false {
		return nil
	}

	v := Task{
		LogID:    logID,
		InsertID: insertID,
		Time:     time.Now(),
		Record:   record,
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
	_, err = taskqueue.Add(c, &task, QUEUE_NAME)
	if err != nil {
		return errors.WrapOr(err)
	}

	return nil
}
