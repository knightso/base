package main

import (
	internal "github.com/knightso/base/errors/internal/unusedroot"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(internal.Analyzer)
}
