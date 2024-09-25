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
	input := `*start
こんにちは[l]世界[p]
12345
aiueo
[image source="test.png"]
[wait time="100"]
*example1
これはexample1です[l]
[jump target="start"]`

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program, err := p.ParseProgram()
	assert.NoError(t, err)

	e := NewEvaluator()
	e.Eval(program)
	{
		events := e.Play("start")
		results := []string{}
		for _, e := range events {
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
			default:
				t.Errorf("未処理のイベントが指定された: %v", event)
			}
		}
		expect := []string{
			"こんにちは",
			"lineEndWait",
			"世界",
			"flush",
			"12345",
			"aiueo",
			"changeBg source=test.png",
			"wait time=100ms",
		}
		assert.Equal(t, expect, results)
	}
	{
		events := e.Play("example1")
		results := []string{}
		for _, e := range events {
			switch event := e.(type) {
			case *worker.MsgEmit:
				results = append(results, event.Body)
			case *worker.LineEndWait:
				results = append(results, "lineEndWait")
			case *worker.Jump:
				results = append(results, fmt.Sprintf("jump target=%s", event.Target))
			default:
				t.Errorf("未処理のイベントが指定された: %v", event)
			}
		}
		expect := []string{
			"これはexample1です",
			"lineEndWait",
			"jump target=start",
		}
		assert.Equal(t, expect, results)
	}
}
