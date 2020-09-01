package main

import (
	"github.com/MakotoNaruse/complexfunc"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(complexfunc.Analyzer) }
