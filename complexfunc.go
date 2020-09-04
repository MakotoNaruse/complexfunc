package complexfunc

import (
	"fmt"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/ssa"
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

func run(pass *analysis.Pass) (interface{}, error) {
	fmt.Println("Calculated by AST Search")
	calcByAST(pass)
	fmt.Println("Calculated by SSA and Control Graph")
	calcBySSA(pass)
	return nil, nil
}

func calcBySSA(pass *analysis.Pass) (interface{}, error) {
	s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	for _, f := range s.SrcFuncs {
		fmt.Println("func name:", f.Name())
		complex := complexity(f)
		fmt.Println("complex:", complex)
		if complex > 10 {
			pass.Reportf(f.Pos(), "function %s is too complicated %d > 10", f.Name(), complex)
		}
	}
	return nil, nil
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
	for _, b := range fn.Blocks {
		edges += len(b.Succs)
	}
	nodes := len(fn.Blocks)
	fmt.Println("n:", nodes, "e:", edges)
	return edges - nodes + 2
}

func calcByAST(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			fmt.Println("func name:", n.Name)
			complex := 1
			complex += calcComplex(n)
			fmt.Println("complex", complex)
		}
	})

	return nil, nil
}

func calcComplex(node ast.Node) int {
	complex := 0
	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause, *ast.CommClause:
			complex++
		case *ast.BinaryExpr:
			if n.Op == token.LAND || n.Op == token.LOR {
				complex++
			}
		}
		return true
	})
	return complex
}
