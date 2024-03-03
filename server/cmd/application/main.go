package main

import (
	"YaGoCalc2/server/internal/calculations/expression"
	"YaGoCalc2/server/internal/orchestrator"
)

func main() {
	orch := orchestrator.NewOrchestrator()
	orch.InitAgent(0)
	go orch.ManageTasks()
	go orch.SendResults()
	expr, err := expression.NewExpression("3+3+3+3+3")
	if err != nil {
		panic(err)
	}
	expr.ParseAST(orch.ExpressionChan, orch.ResultChan)

}
