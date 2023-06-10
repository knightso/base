package internal

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
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			ast.Inspect(funcDecl, func(node ast.Node) bool {
				if ifStmt, ok := node.(*ast.IfStmt); ok {
					if !isRootUsed(pass.TypesInfo, ifStmt.Cond) {
						pass.Report(analysis.Diagnostic{
							Pos:     ifStmt.Pos(),
							Message: "use `errors.Root` to compare with error type",
						})
					}
				}

				if typeAssertExpr, ok := node.(*ast.TypeAssertExpr); ok {
					if !isTypeAssertionWithoutRootUsed(pass.TypesInfo, typeAssertExpr) {
						pass.Report(analysis.Diagnostic{
							Pos:     typeAssertExpr.Pos(),
							Message: "use `errors.Root` before type assertion",
						})
					}
				}

				return true
			})
		}
	}

	return nil, nil
}

// check if the condition of if statement is "errors.Root(err) != ErrHoge"
func isRootUsed(typesInfo *types.Info, cond ast.Expr) bool {
	binaryExpr, ok := cond.(*ast.BinaryExpr)
	if !ok {
		return true
	}

	if !isErrorType(typesInfo, binaryExpr.X) {
		return true
	}
	if callExpr, ok := binaryExpr.X.(*ast.CallExpr); ok {
		// do not allow except errors.Root such as errors.Cause, etc.
		if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := selectorExpr.X.(*ast.Ident); ok {
				if ident.Name == "errors" && selectorExpr.Sel.Name == "Root" {
					return true
				}
			}
		}
	}
	if ident, ok := binaryExpr.Y.(*ast.Ident); ok {
		return ident.Name == "nil"
	}

	return false
}

// check if the type assertion of error is occured
func isTypeAssertionWithoutRootUsed(typesInfo *types.Info, typeAssertExpr *ast.TypeAssertExpr) bool {
	ident, ok := typeAssertExpr.X.(*ast.Ident)
	if !ok {
		return true
	}
	return !isErrorType(typesInfo, ident)
}

func isErrorType(typesInfo *types.Info, expr ast.Expr) bool {
	namedType, ok := typesInfo.TypeOf(expr).(*types.Named)
	if !ok {
		return false
	}
	return namedType.Origin().Obj().Name() == "error"
}
