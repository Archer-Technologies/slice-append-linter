package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "slice-append-linter",
	Doc:      "Checks that .append() calls are not assigned to a new variable.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		// Only assignments should be checked
		(*ast.AssignStmt)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		assignStmt := node.(*ast.AssignStmt)

		if len(assignStmt.Lhs) < 1 {
			return
		}

		// Is this an append() call?
		rhCall, ok := assignStmt.Rhs[0].(*ast.CallExpr)
		if !ok {
			return
		}

		rhCallFunIdent, ok := rhCall.Fun.(*ast.Ident)
		if rhCallFunIdent.Name != "append" {
			return
		}

		// What is being appended / assigned?
		appendSource := rhCall.Args[0]
		appendSourceIdent, ok := appendSource.(*ast.Ident)
		if !ok {
			return
		}

		assignTarget := assignStmt.Lhs[0]
		assignTargetIdent, ok := assignTarget.(*ast.Ident)
		if !ok {
			return
		}

		// Is the assignment to a different variable?
		if appendSourceIdent.Name == assignTargetIdent.Name {
			return
		}

		pass.Reportf(node.Pos(), "slice append() should not assign from source variable %s to different variable %s",
			appendSource, assignTarget)
	})

	return nil, nil
}
