package gae

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"appengine"

	"code.google.com/p/go-uuid/uuid"
	"github.com/go-martini/martini"
	"github.com/knightso/base/gae/bq"
	"github.com/martini-contrib/render"
)

type ExMartini struct {
	*martini.Martini
	martini.Router
}

type MartiniOption struct {
	AdditionalHandlers []martini.Handler
	Log2bq             bool
}

func NewMartini(option MartiniOption) *ExMartini {
	r := martini.NewRouter()
	m := martini.New()

	m.Use(func(c martini.Context, r *http.Request, l *log.Logger) {
		ac := appengine.NewContext(r)
		gaelog := log.New(logWriter{ac}, l.Prefix(), l.Flags())
		c.Map(gaelog)
	})
	m.Use(func(c martini.Context, r *http.Request) {
		ac := appengine.NewContext(r)
		c.Next()
		for l := popLog(); l != ""; l = popLog() {
			ac.Debugf(l)
			
			if option.Log2bq == false {
				continue
			}

			id := uuid.NewUUID()
			uuidString := id.String()
			now := time.Now()
			record := make(map[string]interface{})
			record["id"] = uuidString
			record["date"] = now
			record["log"] = l

			err := bq.SendLog(ac, "debuglog", uuidString, record)
			if err != nil {
				ac.Warningf("%s", err.Error())
				continue
			}
		}
	})
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Use(render.Renderer(render.Options{
		IndentJSON: true,
	}))
	for _, h := range option.AdditionalHandlers {
		m.Use(h)
	}
	m.MapTo(r, (*martini.Route)(nil))
	m.Action(r.Handle)

	return &ExMartini{m, r}
}

type logWriter struct {
	ac appengine.Context
}

func (w logWriter) Write(p []byte) (n int, err error) {
	w.ac.Debugf(string(p))
	return len(p), nil
}

var logs []string
var logMutex sync.Mutex

var LOG_ENABLED bool

func init() {
	logs = make([]string, 0)
}

func Logf(s string, a ...interface{}) {
	if !LOG_ENABLED {
		return
	}

	logMutex.Lock()
	defer logMutex.Unlock()

	logs = append(logs, location()+": "+fmt.Sprintf(s, a...))
}

func popLog() string {
	if !LOG_ENABLED {
		return ""
	}

	logMutex.Lock()
	defer logMutex.Unlock()

	if len(logs) == 0 {
		return ""
	}

	s := logs[0]
	logs = logs[1:]

	return s
}

func location() string {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
		line = -1
	}
	return fmt.Sprintf("%s:%d", file, line)
}
