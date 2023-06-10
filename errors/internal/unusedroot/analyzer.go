package unusedroot

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "unusedroot",
	Doc:  "unusedroot checks whether if errors.Root is used",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	rules := []Rule{
		&RuleRootUsedForTypeAssertion{pass.TypesInfo},
		&RuleRootUsedForCompare{pass.TypesInfo},
	}
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			ast.Inspect(funcDecl, func(node ast.Node) bool {
				for _, rule := range rules {
					if rule.IsTarget(node) {
						if !rule.Check(node) {
							pass.Report(rule.Diagnostic(node))
						}
					}
				}
				return true
			})
		}
	}

	return nil, nil
}

// call check() before call isTarget and only if its result is true
type Rule interface {
	IsTarget(ast.Node) bool
	Check(ast.Node) bool
	Diagnostic(ast.Node) analysis.Diagnostic
}

func isErrorType(typesInfo *types.Info, expr ast.Expr) bool {
	namedType, ok := typesInfo.TypeOf(expr).(*types.Named)
	if !ok {
		return false
	}
	return namedType.Origin().Obj().Name() == "error"
}

func isRootedExpr(expr ast.Expr) bool {
	switch expr := expr.(type) {
	case *ast.CallExpr:
		// do not allow except "errors.Root" such as "errors.Cause", etc.
		if selectorExpr, ok := expr.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := selectorExpr.X.(*ast.Ident); ok {
				if ident.Name == "errors" && selectorExpr.Sel.Name == "Root" {
					return true
				}
			}
		}
		return false
	case *ast.Ident:
		// TODO: "errors.Root"-ed かどうかを確認する
		return false
	default:
		return true
	}
}
