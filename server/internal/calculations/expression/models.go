package expression

import (
	"go/ast"
	"go/token"
)

type Expression struct {
	RawExpression string
	ASTExpr       ast.Expr
	Fset          *token.FileSet
}
