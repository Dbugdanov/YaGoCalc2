package agent

import (
	"YaGoCalc2/server/internal/calculations/expression"
)

type Status int

const (
	ALIVE   Status = iota // Значит что можно поручить задачку
	WORKING               // Значит что живой, но работает

)

type Agent struct {
	Id         int                        // Идентификатор агента
	Status     Status                     // статус агента
	StatusChan chan Status                // Канал для отправки статуса оркестратору
	Tasks      chan expression.Task       // Канал для получения заданий на вычисление
	Report     chan expression.TaskResult // Канал для отправки результатов вычисления
}
