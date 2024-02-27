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

func (expr *Expression) ParseAST() {
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
					var val = constant.BinaryOp(constant.MakeFromLiteral(leftExpr.Value, leftExpr.Kind, 0), x.Op, constant.MakeFromLiteral(rightExpr.Value, rightExpr.Kind, 0))
					fmt.Println("Visitor found an expression. Sending expression to agent. Visitor position: ", x.OpPos, "Agent calculated: ", val, " KIND: ", val.Kind())
					c.Replace(&ast.BasicLit{
						ValuePos: n.End(),
						Kind:     token.FLOAT,
						Value:    val.String(),
					})
				}
			}
		}

		return true
	})
}
