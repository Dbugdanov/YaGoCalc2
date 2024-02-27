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

func (orch *Orchestrator) ListAllAgents() {
	fmt.Println(orch.pool.agents)
}

func (orch *Orchestrator) initAgent(id int) *agent.Agent {
	newAgent, ok := agent.NewAgent(id)
	if ok {
		orch.pool.agents = append(orch.pool.agents, newAgent)

	}
	return newAgent
}

func (orch *Orchestrator) RunAgentPool() {

	var wg sync.WaitGroup

	wg.Add(2)
	// 1. из экземпляра пула создать горутины
	for i := 0; i < 2; i++ {
		go func(i int) {
			agent := orch.initAgent(i)
			agent.Run()
		}(i)

	}
}
