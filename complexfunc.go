package complexfunc

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/ssa"
	"sort"
)

const doc = "complexfunc is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "complexfunc",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		buildssa.Analyzer,
	},
}

var (
	over  int
)

func init() {
	flag.IntVar(&over, "over", 10, "report functions which has complexity > over")
}

func run(pass *analysis.Pass) (interface{}, error) {
	scores := map[token.Pos]score{}

	resultBySsa := calcBySSA(pass, scores)
	result := calcByAST(pass, resultBySsa)
	// sort by pos
	var keys []int
	for k := range result {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, k := range keys {
		pos := token.Pos(k)
		score := scores[pos]
		fmt.Println(score)
		if score.astCmp > over {
			pass.Reportf(pos, "function %s.%s is too complicated %d > 10", score.PkgName, score.FuncName, score.astCmp)
		}
		if score.ssaCmp < score.astCmp {
			pass.Reportf(pos, "function %s.%s has redundant branch", score.PkgName, score.FuncName)
		}
	}
	return nil, nil
}

type score struct {
	PkgName  string
	FuncName string
	astCmp   int
	ssaCmp   int
	Pos      token.Pos
}

func (s score) String() string {
	return fmt.Sprintf("function: %s.%s\nscore by ast: %d\nscore by ssa: %d", s.PkgName, s.FuncName, s.astCmp, s.ssaCmp)
}

func calcBySSA(pass *analysis.Pass, scores map[token.Pos]score) map[token.Pos]score {
	s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	for _, f := range s.SrcFuncs {
		scores[f.Pos()] = score{
			PkgName:  f.Pkg.Pkg.Name(),
			FuncName: f.Name(),
			astCmp:   0,
			ssaCmp:   complexity(f),
			Pos:      f.Pos(),
		}
	}
	return scores
}

func complexity(fn *ssa.Function) int {
	/*
		https://en.wikipedia.org/wiki/Cyclomatic_complexity
		The complexity M for a function is defined as
		M = E − N + 2
		E = the number of edges.
		N = the number of nodes.
	*/
	edges := 0
	returns := 0
	for _, b := range fn.Blocks {
		edges += len(b.Succs)
		for _, instr := range b.Instrs {
			switch instr.(type) {
			case *ssa.Return:
				returns++
			}
		}
	}
	nodes := len(fn.Blocks)
	//fmt.Println("n:", nodes, "e:", edges)
	return edges - nodes + 2 + returns - 1
}

func calcByAST(pass *analysis.Pass, scores map[token.Pos]score) map[token.Pos]score {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			score , ok := scores[n.Name.NamePos]
			if !ok { return }
			score.astCmp = 1 + calcComplex(n)
			scores[n.Name.NamePos] = score
		}
	})

	return scores
}

func calcComplex(node ast.Node) int {
	complex := 0
	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt:
			complex++
		case *ast.CaseClause:
			// len == 0 is default cause
			if len(n.List) > 0 {
				complex++
			}
		case *ast.CommClause:
			// Comm == nil is default cause
			if n.Comm != nil {
				complex++
			}
		case *ast.BinaryExpr:
			if n.Op == token.LAND || n.Op == token.LOR {
				complex++
			}
		}
		return true
	})
	return complex
}
