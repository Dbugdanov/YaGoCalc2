package orchestrator

import (
	"YaGoCalc2/server/internal/agent"
	"YaGoCalc2/server/internal/calculations/expression"
	"fmt"
	"go/constant"
	"sync"
)

// отслеживать количество агентов
// отслеживать состояние каждого агента
// создать нового агента
// удалить определенного агента
// принять выражение и выдать задачу на вычисление свободному агенту с определенным временем

// NewOrchestrator - конструктор для создания нового оркестратора. создаётся 1 раз в application. Должен как и парсер выражений быть демоном, отслеживать событие
// Оркестратор должен уметь отслеживать новые части выражения и отдавать их агентам.
func NewOrchestrator() *Orchestrator {
	ap := AgentPool{
		agents:     make([]*agent.Agent, 0),
		taskChan:   make(chan expression.Task),
		resultChan: make(chan constant.Value),
	}
	return &Orchestrator{
		pool: ap,
	}
}

// ListAllAgents возможно это лишнее
func (orch *Orchestrator) ListAllAgents() {
	fmt.Println(orch.pool.agents)
}

// initAgent инициализировать нового агента. Агент это горутина которая умеет принимать задачу с выражением и отдавать ответ.
func (orch *Orchestrator) initAgent(id int) *agent.Agent {
	newAgent, ok := agent.NewAgent(id)
	if ok {
		orch.pool.agents = append(orch.pool.agents, newAgent)
		orch.pool.resultChan <- newAgent.Run(orch.pool.taskChan)
	}
	return newAgent
}

// RunAgentPool запустить пул агентов. Пока что захардкожено всего 2 воркера (агента).
// пул агентов получает задачи от оркестратора и после того как выполняет отдаёт оркестратору. А оркестратор уже отдаёт обратно в парсер.
func (orch *Orchestrator) RunAgentPool() {

	var wg sync.WaitGroup

	wg.Add(2)
	// 1. из экземпляра пула создать горутины
	for i := 0; i < 2; i++ {
		defer wg.Done()
		go func(i int) {
			orch.initAgent(i)
		}(i)
		wg.Wait()
	}
}
