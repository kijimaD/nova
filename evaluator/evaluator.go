package evaluator

import (
	"fmt"
	"time"

	"github.com/kijimaD/nov/ast"
	"github.com/kijimaD/nov/token"
	"github.com/kijimaD/nov/worker"
)

type Evaluator struct {
	Events   []worker.Event
	LabelMap map[string]*ast.BlockStatement
	errors   []error
}

func NewEvaluator() *Evaluator {
	e := Evaluator{
		LabelMap: make(map[string]*ast.BlockStatement),
		Events:   []worker.Event{},
	}

	return &e
}

func (e *Evaluator) Eval(node ast.Node) worker.Event {
	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node)
	case *ast.ExpressionStatement:
		return e.Eval(node.Expression)
	case *ast.BlockStatement:
		for _, statement := range node.Statements {
			e.Eval(statement)
		}
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
				e.errors = append(e.errors, err)
				return nil
			}
			eve = &worker.Wait{DurationMsec: duration}
		case token.CMD_JUMP:
			eve = &worker.Jump{Target: node.Parameters.Map["target"]}
		}
		e.Events = append(e.Events, eve)
		return eve
	case *ast.TextLiteral:
		m := &worker.MsgEmit{Body: node.Value, DoneChan: make(chan bool, 1)}
		e.Events = append(e.Events, m)
		return m
	case *ast.LabelLiteral:
		e.LabelMap[node.LabelName.String()] = node.Body
		return e.Eval(node.Body)
	case nil:
	default:
		e.errors = append(e.errors, fmt.Errorf("error: 未登録のASTを検知した %#v", node))
		return nil
	}

	return nil
}

// 指定ラベルを再生する
func (e *Evaluator) Play(label string) []worker.Event {
	block, ok := e.LabelMap[label]
	if !ok {
		e.errors = append(e.errors, fmt.Errorf("指定ラベルが存在しない %s", label))
		return e.Events
	}
	e.Events = []worker.Event{}
	e.Eval(block)
	return e.Events
}

func (e *Evaluator) evalProgram(program *ast.Program) worker.Event {
	var result worker.Event

	for _, statement := range program.Statements {
		result = e.Eval(statement)
	}

	return result
}
