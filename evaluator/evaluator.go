package evaluator

import (
	"fmt"
	"log"
	"time"

	"github.com/kijimaD/nov/ast"
	"github.com/kijimaD/nov/token"
	"github.com/kijimaD/nov/worker"
)

type Evaluator struct {
	Events []worker.Event
}

func NewEvaluator(node ast.Node) *Evaluator {
	e := Evaluator{}
	e.Eval(node)

	return &e
}

func (e *Evaluator) Eval(node ast.Node) worker.Event {
	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node)
	case *ast.ExpressionStatement:
		return e.Eval(node.Expression)
	case *ast.CmdLiteral:
		var eve worker.Event
		switch node.FuncName.Value {
		case token.CMD_FLUSH:
			eve = &worker.Flush{}
		case token.CMD_LINE_END_WAIT:
			eve = &worker.LineEndWait{}
		case token.CMD_IMAGE:
			eve = &worker.ChangeBg{Source: node.Parameters.Map["source"]}
		case token.CMD_WAIT:
			duration, err := time.ParseDuration(fmt.Sprintf("%sms", node.Parameters.Map["time"]))
			if err != nil {
				log.Fatal(err)
			}
			eve = &worker.Wait{DurationMsec: duration}
		}
		e.Events = append(e.Events, eve)
		return eve
	case *ast.TextLiteral:
		m := &worker.MsgEmit{Body: node.Value, DoneChan: make(chan bool, 1)}
		e.Events = append(e.Events, m)
		return m
	}

	return nil
}

func (e *Evaluator) evalProgram(program *ast.Program) worker.Event {
	var result worker.Event

	for _, statement := range program.Statements {
		result = e.Eval(statement)
	}

	return result
}
