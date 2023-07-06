package unusedroot

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type RuleRootUsedForCompare struct {
	typesInfo *types.Info
}

// check if the node is "if {X} {op} {Y}" and the type of {X} and {Y} is error
func (r *RuleRootUsedForCompare) IsTarget(node ast.Node) bool {
	if ifStmt, ok := node.(*ast.IfStmt); ok {
		if binaryExpr, ok := ifStmt.Cond.(*ast.BinaryExpr); ok {
			return isErrorType(r.typesInfo, binaryExpr.X)
		}
	}
	return false
}

// check if the binaryExpr is not "if {X} {op} nil" and {X} is "errors.Root({val})"-ed
func (r *RuleRootUsedForCompare) Check(node ast.Node) bool {
	binaryExpr := node.(*ast.IfStmt).Cond.(*ast.BinaryExpr)

	// if the compare is only whether the value is nil, it does not needed to check
	if ident, ok := binaryExpr.Y.(*ast.Ident); ok {
		if ident.Name == "nil" {
			return true
		}
	}
	return isRootedExpr(binaryExpr.X)
}

func (r *RuleRootUsedForCompare) Diagnostic(node ast.Node) analysis.Diagnostic {
	return analysis.Diagnostic{
		Pos:     node.Pos(),
		Message: "use `errors.Root` to compare with error type",
	}
}
