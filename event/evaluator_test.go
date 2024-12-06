package event

import (
	"fmt"
	"testing"

	"github.com/kijimaD/nova/lexer"
	"github.com/kijimaD/nova/parser"

	"github.com/stretchr/testify/assert"
)

// TODO: テストで中身の処理を書きすぎているのを直す
func TestEval(t *testing.T) {
	input := `*start
こんにちは[l]世界[p]
12345[r]
aiueo[r]
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
			case *Newline:
				results = append(results, "newline")
			default:
				t.Errorf("未処理のイベントが指定された: %#v", event)
			}
		}
		expect := []string{
			"こんにちは",
			"lineEndWait",
			"世界",
			"flush",
			"12345",
			"newline",
			"aiueo",
			"newline",
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

func TestEval_変換できる(t *testing.T) {
	input := `*start
hello
world
`

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program, err := p.ParseProgram()
	assert.NoError(t, err)

	e := NewEvaluator()
	e.Eval(program)
	e.Play("start")

	var estr string
	for _, e := range e.Events {
		estr += fmt.Sprintf("%s\n", e.String())
	}

	expect := `<MsgEmit hello>
<MsgEmit world>
`
	assert.Equal(t, expect, estr)
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
			expect: []string{"start", "ch1", "ch2"},
		},
		{
			name: "定義順に取得する",
			input: `*ch2
ch2
*ch1
ch1
*start
start
`,
			expect: []string{"ch2", "ch1", "start"},
		},
		{
			name:   "空文字の場合は空スライスを返す",
			input:  ``,
			expect: []string{},
		},
		{
			name: "重複していた場合は1つにまとめられる",
			input: `*ch1
xxx
*ch1
xxx
*ch1
xxx
`,
			expect: []string{"ch1"},
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
			assert.Equal(t, 0, len(e.errors))
		})
	}
}
