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
				switch ifs.Else.(type) {
				case *ast.IfStmt:
					complex += calcComplex([]ast.Stmt{ifs.Else})
				default:
					complex += 1 + calcComplex([]ast.Stmt{ifs.Else})
				}
			}
		case *ast.ForStmt:
			frs, _ := stmt.(*ast.ForStmt)
			complex += 1 + calcComplex(frs.Body.List)
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
