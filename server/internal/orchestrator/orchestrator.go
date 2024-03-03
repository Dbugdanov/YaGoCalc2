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

// NewOrchestrator - конструктор для создания нового оркестратора. Создаётся 1 раз в application. Должен, как и парсер выражений быть демоном, отслеживать событие
// Оркестратор должен уметь отслеживать новые части выражения и отдавать их агентам.
func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		AgentPool{
			Agents:      make(map[int]*agent.Agent),
			ResultChan:  make(chan expression.TaskResult),
			AgentStatus: make(chan agent.Status),
			mu:          sync.Mutex{},
		},
		make(chan expression.SubExpression),
		make(chan constant.Value),
	}
}

func (orch *Orchestrator) ManageTasks() {
	for subExpression := range orch.ExpressionChan {
		fmt.Println(subExpression, " received")
		task := expression.NewTask(0, subExpression)
		fmt.Println(task, " created")
		for _, _agent := range orch.Pool.Agents { // Ищем свободного агента
			fmt.Println("current pool of agents: ", orch.Pool.Agents)
			if _agent.Status == agent.ALIVE {
				fmt.Println("found alive agent")
				_agent.Tasks <- *task // Отправляем задачу агенту
				fmt.Println("sent task to agent")
				break // Выходим из цикла после распределения задачи
			}
		}
	}
}

func (orch *Orchestrator) SendResults() {
	for result := range orch.Pool.ResultChan {
		fmt.Println("orch received")
		orch.ResultChan <- result.Result
		fmt.Println("orch sent res")
	}
}

// ListAllAgents возможно это лишнее либо создавалось для того, чтоб отдавать фронту фронта
func (orch *Orchestrator) ListAllAgents() {
	fmt.Println(orch.Pool.Agents)
}

// InitAgent инициализировать нового агента. Агент это горутина которая умеет принимать задачу с выражением и отдавать ответ.
func (orch *Orchestrator) InitAgent(id int) {

	orch.Pool.mu.Lock()
	defer orch.Pool.mu.Unlock()
	// TODO: создание агента наверно вывести отсюда куда-нибудь. Обработать ошибку если не ok
	newAgent, ok := agent.NewAgent(id, orch.Pool.ResultChan)
	if ok {
		orch.Pool.Agents[newAgent.Id] = newAgent
		go newAgent.Run()
	}
}
