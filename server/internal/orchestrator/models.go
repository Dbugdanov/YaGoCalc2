package orchestrator

import (
	"YaGoCalc2/server/internal/agent"
	"YaGoCalc2/server/internal/calculations/expression"
	"go/constant"
	"sync"
)

var id int

type Orchestrator struct {
	Pool           AgentPool
	ExpressionChan chan expression.SubExpression // Канал для получения заданий от парсера
	ResultChan     chan constant.Value
}

type AgentPool struct {
	Agents      map[int]*agent.Agent // Мапа, в которой хранится информация об агентах
	ResultChan  chan expression.TaskResult
	AgentStatus chan agent.Status // Канал для отслеживания статуса агентов
	mu          sync.Mutex        // Синхронизация доступа к агентам
}
