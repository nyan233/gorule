package gorule

import (
	"go/ast"
	"go/token"
)

func numberEval(expr ast.Expr, args map[string]interface{}) interface{} {
	binExpr := expr.(*ast.BinaryExpr)
	switch binExpr.Op {
	case token.ADD:
		break
	case token.SUB:
		break
	case token.MUL:
		break
	case token.QUO:
		break
	}
	return nil
}
