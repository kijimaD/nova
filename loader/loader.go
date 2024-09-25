package loader

import (
	"github.com/kijimaD/nov/event"
	"github.com/kijimaD/nov/lexer"
	"github.com/kijimaD/nov/parser"
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
	q := event.NewQueue()
	q.Events = e.Events

	return q, nil
}
