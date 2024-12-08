package event

import (
	"testing"

	"github.com/kijimaD/nova/lexer"
	"github.com/kijimaD/nova/parser"

	"github.com/stretchr/testify/assert"
)

func TestEval_Snapshot(t *testing.T) {
	input := `*start
こんにちは[l]世界[p]
12345[r]
aiueo[r]
[image source="test.png"]
[wait time="100"]
*example1
これはexample1です[l]
[jump target="start"]`
	q := prepareQueue(t, input)
	q.Start()
	{
		q.Play("start")
		expect := []string{
			"<MsgEmit こんにちは>",
			"<LineEndWait>",
			"<MsgEmit 世界>",
			"<Flush>",
			"<MsgEmit 12345>",
			"<Newline>",
			"<MsgEmit aiueo>",
			"<Newline>",
			"<ChangeBg test.png>",
			"<Wait 100ms>",
		}
		assert.Equal(t, expect, q.DumpQueue())
	}
	{
		q.Play("example1")
		expect := []string{
			"<MsgEmit これはexample1です>",
			"<LineEndWait>",
			"<Jump start>",
		}
		assert.Equal(t, expect, q.DumpQueue())
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
