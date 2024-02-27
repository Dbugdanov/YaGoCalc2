package expression

import (
	"go/ast"
	"go/token"
	"time"
)

type taskStatus int

const (
	NEEDTOCALC taskStatus = iota
	DONE
	FAILED
)

// Expression - основная структура, содержит исходное выражение в виде строки, бинарное дерево выражений, Fset
type Expression struct {
	RawExpression string
	ASTExpr       ast.Expr
	Fset          *token.FileSet
}

// SubExpression - структура для отправки подвыражения агенту. Представляет собой бинарное выражение. Отправляется в таком
// виде агенту, он возвращает вычисленное значение.
type SubExpression struct {
	Left  ast.BasicLit
	Op    token.Token
	Right ast.BasicLit
}

type Task struct {
	TaskID     int
	SubExpr    SubExpression
	TimeToExec time.Time
	TaskStatus taskStatus
}
