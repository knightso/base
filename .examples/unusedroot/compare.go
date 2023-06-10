package main

import (
	"log"
	"os"

	"github.com/knightso/base/errors"
)

func main() {
	fileInfo, err := os.Stat("hoge")
	if err != nil {
		if err == os.ErrExist { // NG
			log.Fatal(`the path "hoge" does not exist`, err)
		}
		if errors.Root(err) == os.ErrExist { // OK
			log.Fatal(`the path "hoge" does not exist`, err)
		}
		log.Fatal("unknown error", err)
	}
	_ = fileInfo
}
