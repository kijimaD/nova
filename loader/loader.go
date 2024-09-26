package loader

import (
	"github.com/kijimaD/nova/event"
	"github.com/kijimaD/nova/lexer"
	"github.com/kijimaD/nova/parser"
)

// 依存関係のせいで、適当に配置できない
// スクリプトからキューを初期化する
func NewQueueFromText(text string) (event.Queue, error) {
	l := lexer.NewLexer(text)
	p := parser.NewParser(l)
	program, err := p.ParseProgram()
	if err != nil {
		return event.Queue{}, err
	}
	e := event.NewEvaluator()
	e.Eval(program)
	q := event.NewQueue(e)
	q.Evaluator.Events = e.Events

	return q, nil
}
