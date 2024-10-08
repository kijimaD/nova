package event

import (
	"fmt"
	"testing"

	"github.com/kijimaD/nova/lexer"
	"github.com/kijimaD/nova/parser"

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
		e.Play("start")
		results := []string{}
		for _, e := range e.Events {
			switch event := e.(type) {
			case *MsgEmit:
				results = append(results, event.Body)
			case *Flush:
				results = append(results, "flush")
			case *LineEndWait:
				results = append(results, "lineEndWait")
			case *ChangeBg:
				results = append(results, fmt.Sprintf("changeBg source=%s", event.Source))
			case *Wait:
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
		e.Play("example1")
		results := []string{}
		for _, e := range e.Events {
			switch event := e.(type) {
			case *MsgEmit:
				results = append(results, event.Body)
			case *LineEndWait:
				results = append(results, "lineEndWait")
			case *Jump:
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

func TestList(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect []string
	}{
		{
			name: "取得できる",
			input: `*start
start
*ch1
ch1
*ch2
ch2
`,
			expect: []string{"ch1", "ch2", "start"},
		},
		{
			name: "ソートして取得する",
			input: `*ch2
ch2
*ch1
ch1
*start
start
`,
			expect: []string{"ch1", "ch2", "start"},
		},
		{
			name:   "空文字の場合は空スライスを返す",
			input:  ``,
			expect: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			p := parser.NewParser(l)
			program, err := p.ParseProgram()
			assert.NoError(t, err)
			e := NewEvaluator()
			e.Eval(program)
			result := e.Labels()

			assert.Equal(t, tt.expect, result)
		})
	}
}
