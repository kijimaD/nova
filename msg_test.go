package msg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMsg(t *testing.T) {
	q := Queue{}
	q.events = append(q.events, &msgEmit{
		body: []rune("こんにちは"),
	})
}

func TestEmit(t *testing.T) {
	q := Queue{}
	q.events = append(q.events, &msgEmit{
		body:   []rune("東京"),
		status: TaskNotRunning,
	})
	q.events = append(q.events, &flush{})
	q.events = append(q.events, &msgEmit{
		body:   []rune("京都"),
		status: TaskNotRunning,
	})

	q.Pop()
	assert.Equal(t, "", q.Display())
	time.Sleep(80 * time.Millisecond)
	assert.Equal(t, "東京", q.Display())
}

// func TestWait(t *testing.T) {
// 	q := Queue{}
// 	q.events = append(q.events, &msgEmit{
// 		body: []rune("東京"),
// 	})
// 	q.events = append(q.events, &flush{})
// 	q.events = append(q.events, &msgEmit{
// 		body: []rune("京都"),
// 	})
// }

// func TestBuilder(t *testing.T) {
// 	input := `こんにちは...[p]
// 今日はいかがですか`
// 	l := NewLexer(input)
// 	p := NewParser(l)
// 	program := p.ParseProgram()
// 	e := Evaluator{}
// 	e.Eval(program)

// 	q := NewQueue(e.Events)
// 	assert.Equal(t, "今日はいかがですか", q.buf)
// }

// // 改行を自動挿入できる
// func TestNewLine(t *testing.T) {
// 	input := `こんにちは[p]
// ああああああああああああああああああああ`
// 	l := NewLexer(input)
// 	p := NewParser(l)
// 	program := p.ParseProgram()
// 	e := Evaluator{}
// 	e.Eval(program)

// 	q := NewQueue(e.Events)
// 	assert.Equal(t, "ああああああああああああああ\nあああああ", q.buf)
// }

// // 意図的な改行で自動改行カウントをリセットする
// func TestNewLineResetCount(t *testing.T) {
// 	input := `こんにちは[p]
// ああああああああああ
// ああああああああああ`
// 	l := NewLexer(input)
// 	p := NewParser(l)
// 	program := p.ParseProgram()
// 	e := Evaluator{}
// 	e.Eval(program)

// 	q := NewQueue(e.Events)
// 	assert.Equal(t, "こんにちは", q.buf)
// 	assert.Equal(t, "ああああああああああ\nああああああああああ", q.buf)
// }
