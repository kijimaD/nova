package msg

import (
	"testing"
	"time"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestMsg(t *testing.T) {
	q := Queue{}
	q.events = append(q.events, &msgEmit{
		body: "こんにちは",
	})
}

// func TestPop(t *testing.T) {
// 	q := Queue{}
// 	q.events = append(q.events, &flush{})
// 	q.events = append(q.events, &notImplement{})
// 	q.events = append(q.events, &flush{})
// 	q.events = append(q.events, &notImplement{})

// 	assert.Equal(t, &flush{}, q.Pop())
// 	assert.Equal(t, &notImplement{}, q.Pop())
// 	assert.Equal(t, &flush{}, q.Pop())
// 	assert.Equal(t, &notImplement{}, q.Pop())
// }

// いちいちsleepしないといけないのをどうにかしたい
// テストの場合は同期的にやりたい。終了を通知してくるまで待つとかできないか?
func TestMsgEmit(t *testing.T) {
	q := NewQueue()
	q.events = append(q.events, &msgEmit{
		body: "東京1東京2東京3東京4東京5東京6東京7東京8東京9東京10東京11東京12",
	})
	q.events = append(q.events, &flush{})
	q.Start()

	assert.Equal(t, "", q.Display())
	q.Pop()
	assert.Equal(t, "", q.Display())
	time.Sleep(50 * time.Millisecond)
	assert.True(t, utf8.RuneCountInString(q.Display()) > 1)
	assert.True(t, utf8.RuneCountInString(q.Display()) < 10)
	q.Skip()
	q.Wait()
	assert.Equal(t, "東京1東京2東京3東京4東京5東京6東京7東京8東京9東京10東京11東京12", q.Display())
	q.Pop()
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, "", q.Display())
}

func TestRun_RunがPopとSkipを使い分ける(t *testing.T) {
	q := NewQueue()
	q.events = append(q.events, &msgEmit{
		body: "東京1東京2東京3東京4東京5東京6東京7東京8東京9東京10東京11東京12",
	})
	q.events = append(q.events, &flush{})
	q.Start()

	q.Pop()
	// FIXME: 最初はPopでないといけない...
	// q.Run() // pop
	time.Sleep(20 * time.Millisecond)
	q.Run() // skip
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, "東京1東京2東京3東京4東京5東京6東京7東京8東京9東京10東京11東京12", q.Display())
	q.Run() // pop
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, "", q.Display())
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
