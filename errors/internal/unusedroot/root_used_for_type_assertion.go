package unusedroot

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type RuleRootUsedForTypeAssertion struct {
	typesInfo *types.Info
}

// check if the node is typeAssertExpr and the type of val is error
func (r *RuleRootUsedForTypeAssertion) IsTarget(node ast.Node) bool {
	if typeAssertExpr, ok := node.(*ast.TypeAssertExpr); ok {
		return isErrorType(r.typesInfo, typeAssertExpr.X)
	}
	return false
}

// check if the typeAssertExpr.X is "errors.Root"-ed
func (r *RuleRootUsedForTypeAssertion) Check(node ast.Node) bool {
	typeAssertExpr := node.(*ast.TypeAssertExpr)
	return isRootedExpr(typeAssertExpr.X)
}

func (r *RuleRootUsedForTypeAssertion) Diagnostic(node ast.Node) analysis.Diagnostic {
	return analysis.Diagnostic{
		Pos:     node.Pos(),
		Message: "use `errors.Root` before type assertion",
	}
}
