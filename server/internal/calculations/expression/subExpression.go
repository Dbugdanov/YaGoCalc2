package expression

import (
	"go/ast"
	"go/token"
)

func newSubExpression(
	left ast.BasicLit,
	op token.Token,
	right ast.BasicLit,
) *SubExpression {
	return &SubExpression{
		Left:  left,
		Op:    op,
		Right: right,
	}
}

// sendSubExpression - отправляет подвыражение агенту с временем, которое выставлено у пользователя исходя из сессии
//func (subExpr *SubExpression) sendSubExpression() bool {
//
//}
