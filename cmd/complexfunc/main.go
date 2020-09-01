package main

import (
	"golang.org/x/tools/go/analysis/unitchecker"
	"practice/complexfunc"
)

func main() { unitchecker.Main(complexfunc.Analyzer) }
