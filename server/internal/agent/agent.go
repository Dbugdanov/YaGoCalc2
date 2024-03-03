package agent

import (
	"YaGoCalc2/server/internal/calculations/expression"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"time"
)

//TODO: как общаются агент и оркестратор?? внедрить реализацию handshake'а и дальнейшего для общения

// NewAgent - Конструктор для создания концептуально агента. Он будет иметь методы для вычисления выражений, а также
// иметь возможность общаться с менеджером агентов. С БД агент никак не взаимодействует, у него есть только 2 задачи -
// вычислять выражение с заданным таймаутом (таймаут и выражение отправляет менеджер), прекратить вычисление
// (если менеджер так сказал - ну что поделать, перестаёт вычислять, отправляет менеджеру сигнал что "я ALIVE"), и
// каждые 5 секунд отправлять свой статус для мониторинга менеджером
func NewAgent(id int, reportChan chan expression.TaskResult) (*Agent, bool) {
	return &Agent{
		Id:         id,
		Status:     ALIVE,
		StatusChan: make(chan Status),
		Tasks:      make(chan expression.Task),
		Report:     reportChan,
	}, true
}

// Run - что-то вроде запуска агента. Он должен бесконечно слушать канал с тасками, и, как приходит таск, считать его, после чего отдавать ответ в канал для ответа
func (a *Agent) Run() {
	go a.sendHeartBeat() // отправляем статус агента в отдельной горутине
	for {
		select {
		case task := <-a.Tasks: // Когда прилетает таска
			fmt.Println("agent ", a.Id, "task catch")
			a.Status = WORKING // Меняем статус агента на "работаю"
			left := task.SubExpr.Left
			op := task.SubExpr.Op
			right := task.SubExpr.Right
			time.Sleep(task.SubExpr.TimeToExec) // И считаем выражение с заданным таймаутом.
			out := expression.TaskResult{
				AgentID: a.Id,
				Result:  a.calculateTask(left, op, right)}
			fmt.Println("out: /// ", out)
			a.Report <- out // Отдаём результат в соответствующий канал
			fmt.Println("agent send out")
			task.TaskStatus = expression.DONE
			a.Status = ALIVE
		}
	}
}

func (a *Agent) sendHeartBeat() {
	timer := time.NewTicker(5 * time.Second) // Создаём таймер, который срабатывает каждые 5 секунд
	defer timer.Stop()                       // Убедимся, что таймер будет остановлен при выходе из функции

	for {
		select {
		case <-timer.C:
			fmt.Println(a.Status)
			a.StatusChan <- a.Status // Отправляем текущее состояние агента
		}
	}
}

// calculateTask вычисляет выражение и отдаёт
func (a *Agent) calculateTask(left ast.BasicLit, op token.Token, right ast.BasicLit) constant.Value {
	fmt.Println("calculating...")
	return constant.BinaryOp(constant.MakeFromLiteral(left.Value, left.Kind, 0), op, constant.MakeFromLiteral(right.Value, right.Kind, 0))
}
