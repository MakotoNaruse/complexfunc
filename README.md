# complexfunc
[![Go Report Card](https://goreportcard.com/badge/github.com/MakotoNaruse/complexfunc)](https://goreportcard.com/report/github.com/MakotoNaruse/complexfunc)
[![Go test](https://github.com/MakotoNaruse/complexfunc/workflows/Go%20test/badge.svg?branch=master)](https://github.com/MakotoNaruse/complexfunc/actions)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

`complexfunc` finds too complicated function.

Reported if cyclomatic complexity calculated by SSA exceeds n (default 10).

We also calculate cyclomatic complexity by AST and report if cyclo by SSA < AST

In that case, there are redundant branches and SSA removed them.

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

func f2() int { // want "function a.f2 is too complicated 12 > 10"
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
func f3(a int) int { // want "function a.f3 has redundant branch"
	if false {
	}
	return a + 1
}
```
```console
function: a.f1
score by ast: 4
score by ssa: 4
function: a.f2
score by ast: 12
score by ssa: 12
function: a.f3
score by ast: 2
score by ssa: 1
./a.go:18:1: function a.f1 is too complicated 12 > 10
./a.go:46:1: function a.f3 has redundant branch
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
$ go vet -vettool=`which complecfunc` [flag] pkgname
Flags:
      --over N   report if complexity > N (default 10)
```
