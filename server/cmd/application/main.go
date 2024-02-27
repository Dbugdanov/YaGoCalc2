package main

import (
	"YaGoCalc2/server/internal/calculations/expression"
	"YaGoCalc2/server/internal/orchestrator"
)

func main() {
	expr, err := expression.NewExpression("1+2+3+4*5")
	if err != nil {
		panic(err)
	}
	expr.ParseAST()
	orch := orchestrator.NewOrchestrator()
	orch.RunAgentPool()
	orch.ListAllAgents()

}
