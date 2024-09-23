package evaluator

import (
	"fmt"
	"testing"

	"github.com/kijimaD/nov/lexer"
	"github.com/kijimaD/nov/parser"
	"github.com/kijimaD/nov/worker"

	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	input := `こんにちは[l]世界[p]
←12345
←aiueo
[image source="test.png"]
[wait time="100"]`

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	e := Evaluator{}
	e.Eval(program)
	results := []string{}
	for _, e := range e.Events {
		switch event := e.(type) {
		case *worker.MsgEmit:
			results = append(results, event.Body)
		case *worker.Flush:
			results = append(results, "flush")
		case *worker.LineEndWait:
			results = append(results, "lineEndWait")
		case *worker.ChangeBg:
			results = append(results, fmt.Sprintf("changeBg source=%s", event.Source))
		case *worker.Wait:
			results = append(results, fmt.Sprintf("wait time=%s", event.DurationMsec))
		}
	}
	expect := []string{
		"こんにちは",
		"lineEndWait",
		"世界",
		"flush",
		"←12345",
		"←aiueo",
		"changeBg source=test.png",
		"wait time=100ms",
	}
	assert.Equal(t, expect, results)
}
