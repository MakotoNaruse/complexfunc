package complexfunc

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "complexfunc is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "complexfunc",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			fmt.Println("func name:", n.Name.Name)
			complex := 1
			complex += calcComplex(n.Body.List)
			fmt.Println("complex", complex)
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
			fmt.Println("if cause", ifs.Pos())
			complex += 1 + calcComplex(ifs.Body.List)
			if ifs.Else != nil {
				switch ifs.Else.(type){
				case *ast.IfStmt:
					complex += calcComplex([]ast.Stmt{ifs.Else})
				default:
					fmt.Println("else cause", ifs.Else.Pos())
					complex += 1 + calcComplex([]ast.Stmt{ifs.Else})
				}
			}
		case *ast.ForStmt:
			frs, _ := stmt.(*ast.ForStmt)
			fmt.Println("for cause", frs.Pos())
			complex += 1 + calcComplex(frs.Body.List)
		case *ast.SwitchStmt:
			sws, _ := stmt.(*ast.SwitchStmt)
			fmt.Println("switch cause", sws.Pos())
			complex += calcComplex(sws.Body.List)
		case *ast.CaseClause:
			cas, _ := stmt.(*ast.CaseClause)
			fmt.Println("case cause", cas.Pos())
			complex += 1 + calcComplex(cas.Body)
		}
	}
	return complex
}

