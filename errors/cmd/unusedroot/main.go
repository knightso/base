package main

import (
	"github.com/knightso/base/errors/internal/unusedroot"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(unusedroot.Analyzer)
}
