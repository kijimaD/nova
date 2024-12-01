package event

import (
	"fmt"
	"sort"
	"time"

	"github.com/kijimaD/nova/ast"
	"github.com/kijimaD/nova/token"
)

type Evaluator struct {
	Events          []Event
	CurrentEventIdx int
	LabelMap        map[string]*ast.BlockStatement
	CurrentLabel    string
	errors          []error
}

func NewEvaluator() *Evaluator {
	e := Evaluator{
		Events:          []Event{},
		CurrentEventIdx: 0,
		LabelMap:        make(map[string]*ast.BlockStatement),
		errors:          []error{},
	}

	return &e
}

func (e *Evaluator) Eval(node ast.Node) Event {
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
		var eve Event
		switch node.FuncName.Value {
		case token.CMD_FLUSH:
			eve = &Flush{}
		case token.CMD_LINE_END_WAIT:
			eve = &LineEndWait{}
		case token.CMD_IMAGE:
			eve = &ChangeBg{Source: node.Parameters.Map["source"]}
		case token.CMD_WAIT:
			duration, err := time.ParseDuration(fmt.Sprintf("%sms", node.Parameters.Map["time"]))
			if err != nil {
				e.errors = append(e.errors, err)
				return nil
			}
			eve = &Wait{DurationMsec: duration}
		case token.CMD_JUMP:
			eve = &Jump{Target: node.Parameters.Map["target"]}
		}
		e.Events = append(e.Events, eve)
		return eve
	case *ast.TextLiteral:
		m := &MsgEmit{Body: node.Value, DoneChan: make(chan bool, 1)}
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

// 指定ラベルの内容でEventsを更新する
func (e *Evaluator) Play(label string) {
	e.CurrentLabel = label
	block, ok := e.LabelMap[label]
	if !ok {
		e.errors = append(e.errors, fmt.Errorf("指定ラベルが存在しない %s", label))
		return
	}
	e.CurrentEventIdx = 0
	e.Events = []Event{} // 初期化
	e.Eval(block)
}

func (e *Evaluator) Labels() []string {
	keys := []string{}
	for k, _ := range e.LabelMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func (e *Evaluator) evalProgram(program *ast.Program) Event {
	var result Event

	for _, statement := range program.Statements {
		result = e.Eval(statement)
	}

	return result
}
