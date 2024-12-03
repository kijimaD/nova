package event

import (
	"fmt"
	"time"

	"github.com/kijimaD/nova/ast"
	"github.com/kijimaD/nova/token"
)

type Evaluator struct {
	Events          []Event
	CurrentEventIdx int
	LabelMaster     LabelMaster
	CurrentLabel    string
	errors          []error
}

type LabelMaster struct {
	Labels     []Label
	LabelIndex map[string]int
}

func (master *LabelMaster) GetLabel(key string) (Label, error) {
	idx := master.LabelIndex[key]
	if idx < 0 || len(master.Labels)-1 < idx {
		return Label{}, fmt.Errorf("keyが存在しない")
	}

	return master.Labels[idx], nil
}

type Label struct {
	Name string
	Body *ast.BlockStatement
}

func NewEvaluator() *Evaluator {
	e := Evaluator{
		Events:          []Event{},
		CurrentEventIdx: 0,
		LabelMaster:     LabelMaster{Labels: []Label{}, LabelIndex: map[string]int{}},
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
		label := Label{
			Name: node.LabelName.String(),
			Body: node.Body,
		}
		_, exists := e.LabelMaster.LabelIndex[label.Name]
		if !exists {
			e.LabelMaster.Labels = append(e.LabelMaster.Labels, label)
			e.LabelMaster.LabelIndex[label.Name] = len(e.LabelMaster.Labels) - 1
		}

		return e.Eval(node.Body)
	case nil:
	default:
		e.errors = append(e.errors, fmt.Errorf("error: 未登録のASTを検知した %#v", node))
		return nil
	}

	return nil
}

// 指定ラベルの内容でEventsを更新する
func (e *Evaluator) Play(key string) {
	e.CurrentLabel = key
	e.CurrentEventIdx = 0

	label, err := e.LabelMaster.GetLabel(key)
	if err != nil {
		e.errors = append(e.errors, fmt.Errorf("指定ラベルが存在しない %s", label))

		return
	}

	e.Events = []Event{} // 初期化
	e.Eval(label.Body)
}

func (e *Evaluator) Labels() []string {
	names := []string{}
	for _, label := range e.LabelMaster.Labels {
		names = append(names, label.Name)
	}

	return names
}

func (e *Evaluator) evalProgram(program *ast.Program) Event {
	var result Event

	for _, statement := range program.Statements {
		result = e.Eval(statement)
	}

	return result
}
