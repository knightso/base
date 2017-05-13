package bq

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/taskqueue"

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

func SendLog(c context.Context, logID, insertID string, record interface{}) error {
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

func DatastoreKeyToPath(key *datastore.Key) string {
	if key == nil {
		return ""
	}
	var buf bytes.Buffer
	key2buf(key, &buf)
	return buf.String()
}

func key2buf(key *datastore.Key, buf *bytes.Buffer) {
	if key.Parent() != nil {
		key2buf(key.Parent(), buf)
	}
	if buf.Len() > 0 {
		buf.WriteString(", ")
	}
	buf.WriteString("\"")
	buf.WriteString(key.Kind())
	buf.WriteString("\", ")
	if key.IntID() != 0 {
		buf.WriteString(strconv.FormatInt(key.IntID(), 10))
	} else {
		buf.WriteString("\"")
		buf.WriteString(key.StringID())
		buf.WriteString("\"")
	}
}
