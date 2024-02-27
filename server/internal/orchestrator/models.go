package orchestrator

import (
	"YaGoCalc2/server/internal/agent"
	"YaGoCalc2/server/internal/calculations/expression"
	"go/constant"
)

var id int

type Orchestrator struct {
	pool AgentPool
}

type AgentPool struct {
	agents     []*agent.Agent
	taskChan   chan expression.Task
	resultChan chan constant.Value
}
