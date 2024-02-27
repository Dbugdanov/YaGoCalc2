package agent

import (
	"YaGoCalc2/server/internal/calculations/expression"
	"go/ast"
	"go/constant"
	"go/token"
)

//TODO: как общаются агент и оркестратор?? внедрить реализацию handshake'а и дальнейшего для общения

// NewAgent - Конструктор для создания концептуально агента. Он будет иметь методы для вычисления выражений, а также
// иметь возможность общаться с менеджером агентов. С БД агент никак не взаимодействует, у него есть только 2 задачи -
// вычислять выражение с заданным таймаутом (таймаут и выражение отправляет менеджер), прекратить вычисление
// (если менеджер так сказал - ну что поделать, перестаёт вычислять, отпрвляет менеджеру сигнал что "я ALIVE"), и
// каждые 5 секунд отправлять свой статус для мониторинга менеджером
func NewAgent(id int) (*Agent, bool) {
	return &Agent{
		Id:     id,
		Status: ALIVE,
	}, true
}

func (a *Agent) sendHeartBeat() {
	//TODO: отправить хартбит во время Run
}

// Run - что-то вроде запуска агента. Он должен бесконечно слушать канал с тасками, и, как приходит таск, считать его, после чего отдавать ответ в канал для ответа
func (a *Agent) Run(taskCh <-chan expression.Task) constant.Value {
	var out constant.Value
	for task := range taskCh {
		left := task.SubExpr.Left
		op := task.SubExpr.Op
		right := task.SubExpr.Right
		//TODO: timeToExec
		out = a.calculateTask(left, op, right)
	}
	return out
}

// calculateTask вычисляет выражение и отдаёт
func (a *Agent) calculateTask(left ast.BasicLit, op token.Token, right ast.BasicLit) constant.Value {
	return constant.BinaryOp(constant.MakeFromLiteral(left.Value, left.Kind, 0), op, constant.MakeFromLiteral(right.Value, right.Kind, 0))
}
