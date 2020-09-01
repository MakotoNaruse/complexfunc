package complexfunc

import (
	"fmt"
	"go/ast"
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
		complex := complexity(f)
		fmt.Println("func name:",f)
		fmt.Println("complex:", complex)
	}
	return nil, nil
}
func complexity(fn *ssa.Function) int {
	// https://en.wikipedia.org/wiki/Cyclomatic_complexity
	// The complexity M for a function is defined as
	// M = E − N + 2
	// where
	//
	// E = the number of edges of the graph.
	// N = the number of nodes of the graph.
	edges := 0
	for _, b := range fn.Blocks {
		edges += len(b.Succs)
	}
	return edges - len(fn.Blocks) + 2
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
			complex += calcComplex(n.Body.List)
			fmt.Println("complex", complex)
			if complex > 10 {
				pass.Reportf(n.Pos(), "function %s is too complicated %d > 10", n.Name.Name, complex)
			}
		}
	})

	return nil, nil
}

func calcComplex(stmts []ast.Stmt) int {
	complex := 0
	for _, stmt := range stmts {
		switch stmt.(type) {
		case *ast.IfStmt:
			ifs, _ := stmt.(*ast.IfStmt)
			complex += 1 + calcComplex(ifs.Body.List)
			if ifs.Else != nil {
				// else if can in ifs.Else
				complex += calcComplex([]ast.Stmt{ifs.Else})
			}
		case *ast.ForStmt:
			frs, _ := stmt.(*ast.ForStmt)
			complex += 1 + calcComplex(frs.Body.List)
		case *ast.RangeStmt:
			rgs, _ := stmt.(*ast.RangeStmt)
			complex += 1 + calcComplex(rgs.Body.List)
		case *ast.SwitchStmt:
			sws, _ := stmt.(*ast.SwitchStmt)
			complex += calcComplex(sws.Body.List)
		case *ast.CaseClause:
			cas, _ := stmt.(*ast.CaseClause)
			complex += 1 + calcComplex(cas.Body)
		}
	}
	return complex
}
