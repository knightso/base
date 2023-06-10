package main

import (
	"log"

	"github.com/knightso/base/errors"
	"google.golang.org/appengine"
)

func main() {
	var err error
	if multiError, ok := err.(appengine.MultiError); ok { // NG
		log.Fatal(multiError)
	}
	if multiError, ok := errors.Root(err).(appengine.MultiError); ok { // OK
		log.Fatal(multiError)
	}
	root := errors.Root(err)
	if multiError, ok := root.(appengine.MultiError); ok { // NG
		log.Fatal(multiError)
	}
}
