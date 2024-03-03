package expression

import (
	"go/ast"
	"go/constant"
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
	Left       ast.BasicLit
	Op         token.Token
	Right      ast.BasicLit
	TimeToExec time.Duration
	Result     chan TaskResult
}

// Task - задание, которое отправляется вычислителю (агенту) на выполнение. Содержит в себе подвыражение, а так же дополнительные поля, такие как идентификатор задания, время выполнения выражения, и статус задания.
type Task struct {
	TaskID     int
	SubExpr    SubExpression
	TaskStatus taskStatus
}

// TaskResult - То, что возвращает агент
type TaskResult struct {
	AgentID int
	Result  constant.Value
}
