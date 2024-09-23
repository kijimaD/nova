package loader

import (
	"msg/evaluator"
	"msg/lexer"
	"msg/parser"
	"msg/worker"
)

// 依存関係のせいで、適当に配置できない
// スクリプトからキューを初期化する
func NewQueueFromText(text string) worker.Queue {
	l := lexer.NewLexer(text)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	e := evaluator.NewEvaluator(program)
	q := worker.NewQueue()
	q.Events = e.Events

	return q
}
