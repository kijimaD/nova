package loader

import (
	"github.com/kijimaD/nov/evaluator"
	"github.com/kijimaD/nov/lexer"
	"github.com/kijimaD/nov/parser"
	"github.com/kijimaD/nov/worker"
)

// 依存関係のせいで、適当に配置できない
// スクリプトからキューを初期化する
func NewQueueFromText(text string) worker.Queue {
	l := lexer.NewLexer(text)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	e := evaluator.NewEvaluator()
	e.Eval(program)
	q := worker.NewQueue()
	q.Events = e.Events

	return q
}
