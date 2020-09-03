# complexfunc
[![Go Report Card](https://goreportcard.com/badge/github.com/MakotoNaruse/complexfunc)](https://goreportcard.com/report/github.com/MakotoNaruse/complexfunc)
[![Go test](https://github.com/MakotoNaruse/complexfunc/workflows/Go%20test/badge.svg?branch=master)](https://github.com/MakotoNaruse/complexfunc/actions)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

`complexfunc` finds too complicated function.

Reported if cyclomatic complexity calculated by SSA exceeds 10.

## example
```go
package a

func f1() int {
	num := 10
	ans := 0
	for i := 0; i < num; i++ {
		if i%2 == 0 {
			ans++
		} else if i%3 == 0 {
			ans += 2
		} else {
			ans -= 3
		}
	}
	return ans
}

func f2() int { // want "function f2 is too complicated 12 > 10"
	n := 10
	switch n {
	case 1:
		n++
	case 2:
		n++
	case 3:
		n++
	case 4:
		n++
	case 5:
		n++
	case 6:
		n++
	case 7:
		n++
	case 8:
		n++
	case 9:
		if n % 10 == 0 {
			n++
		}
	case 10:
		n++
	}
	return n
}
```
```console
Calculated by AST Search
func name: f1
complex 4
func name: f2
complex 12
Calculated by SSA and Control Graph
func name: f1
complex: 4
func name: f2
complex: 12
./a.go:18:1: function f1 is too complicated 12 > 10
```

## how to calculate?
### By AST
Counting nodes of `if` `for (range)` `case (in switch)` `&&` `||`.
### By SSA
Define G as the control graph obtained from SSA, 
cyclomatic complexity C is calculated as follows.
```
C = E − N + 2
E : the number of edges of G.
N : the number of nodes of G.
```

## Install

```sh
$ go get github.com/MakotoNaruse/complexfunc/cmd/complexfunc
```

## Usage

```sh
$ go vet -vettool=`which complecfunc` pkgname
```
