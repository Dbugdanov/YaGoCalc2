package expression

import (
	"YaGoCalc2/server/internal/calculations/tokenizer"
	"errors"
	"fmt"
	"go/ast"
	"go/constant"
	"go/parser"
	"go/scanner"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
)

// NewExpression создаёт новое выражение, является функцией-конструктором для структуры Expression. Проверяет, допустимо
// ли для приложения математическое выражение, если нет - отдаёт ошибку
// TODO: принимать новое выражение от веб сервера
func NewExpression(expr string) (*Expression, error) {
	expression := &Expression{
		RawExpression: expr,
		Fset:          token.NewFileSet(),
	}

	err := expression.checkExpressionAllowed()
	if err != nil {
		return nil, err
	}
	astExpr, err := parser.ParseExpr(expr)
	if err != nil {
		return nil, err
	}
	expression.ASTExpr = astExpr
	return expression, nil
}

func newSubExpression(
	left ast.BasicLit,
	op token.Token,
	right ast.BasicLit,
) *SubExpression {
	return &SubExpression{
		Left:   left,
		Op:     op,
		Right:  right,
		Result: make(chan TaskResult),
	}
}

// sendSubExpression - отправляет подвыражение оркестратору с временем, которое выставлено у пользователя исходя из сессии
func (subExpr *SubExpression) sendSubExpression(expressionChan chan<- SubExpression, resultChan chan constant.Value) constant.Value {
	expressionChan <- *subExpr
	fmt.Println(*subExpr, "sent to chan ", expressionChan)
	var out constant.Value
	for result := range resultChan {
		fmt.Println("received result from orch")
		out = result
		break
	}
	fmt.Println(out)
	return out
}

func NewTask(id int, expression SubExpression) *Task {
	return &Task{
		TaskID:     id,
		SubExpr:    expression,
		TaskStatus: NEEDTOCALC,
	}
}

func (expr *Expression) isTokenAllowed(tok token.Token) bool {
	for _, allowedToken := range tokenizer.AllowedTokens {
		if allowedToken == tok {
			return true
		}
	}
	return false
}

func (expr *Expression) checkExpressionAllowed() error {
	src := []byte(expr.RawExpression)
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, 0)

	for {
		_, tok, _ := s.Scan()
		if tok == token.EOF || tok == token.SEMICOLON {
			break
		}
		if !expr.isTokenAllowed(tok) {
			return errors.New(fmt.Sprintf("Found forbidden token %s in expression %s\nAllowed tokens: %v", tok.String(), expr.RawExpression, tokenizer.AllowedTokens))
		}

	}
	return nil
}

func (expr *Expression) ParseAST(expressionChan chan<- SubExpression, resultChan chan constant.Value) {
	astutil.Apply(expr.ASTExpr, nil, func(c *astutil.Cursor) bool {
		ast.Print(expr.Fset, expr.ASTExpr)
		n := c.Node()
		switch x := n.(type) {
		case *ast.BinaryExpr:
			leftExpr, okLeft := x.X.(*ast.BasicLit)
			if okLeft {
				rightExpr, okRight := x.Y.(*ast.BasicLit)
				if okRight {
					//TODO
					//Добавить логирование. Добавить функционал по отправке подвыражения в менеджер агентов

					// TODO: парсер отдаёт подвыражение оркестратору, ждёт когда вычислится выражение
					result := newSubExpression(*leftExpr, x.Op, *rightExpr).sendSubExpression(expressionChan, resultChan)

					fmt.Println("Visitor found an expression. Sending expression to agent. Visitor position: ", x.OpPos, "Agent calculated: ", result, " KIND: ", result.Kind())
					c.Replace(&ast.BasicLit{
						ValuePos: n.End(),
						Kind:     token.FLOAT,
						Value:    result.String(),
					})
				}
			}
		}

		return true
	})
}
