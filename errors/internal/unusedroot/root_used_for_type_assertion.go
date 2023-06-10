package unusedroot

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type RuleRootUsedForTypeAssertion struct {
	typesInfo *types.Info
}

// check if the type assertion of error is occured
func (r *RuleRootUsedForTypeAssertion) IsTarget(node ast.Node) bool {
	if typeAssertExpr, ok := node.(*ast.TypeAssertExpr); ok {
		return isErrorType(r.typesInfo, typeAssertExpr.X)
	}
	return false
}

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
