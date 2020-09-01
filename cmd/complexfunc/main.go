package main

import (
	"practice/complexfunc"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(complexfunc.Analyzer) }

