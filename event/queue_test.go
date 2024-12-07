package event

import (
	"testing"
	"time"
	"unicode/utf8"

	"github.com/kijimaD/nova/lexer"
	"github.com/kijimaD/nova/parser"
	"github.com/kijimaD/nova/utils"

	"github.com/stretchr/testify/assert"
)

func TestMsgEmit_Skipできる(t *testing.T) {
	evaluator := Evaluator{}
	q := NewQueue(&evaluator)
	q.Evaluator.Events = append(q.Evaluator.Events, utils.GetPtr(NewMsgEmit("東京1東京2東京3東京4東京5東京6東京7東京8東京9東京10東京11東京12")))
	q.Evaluator.Events = append(q.Evaluator.Events, &Flush{})
	q.Evaluator.Events = append(q.Evaluator.Events, utils.GetPtr(NewMsgEmit("last")))
	q.Start()

	assert.Equal(t, "", q.Display())
	time.Sleep(30 * time.Millisecond)
	assert.True(t, utf8.RuneCountInString(q.Display()) > 1)
	assert.True(t, utf8.RuneCountInString(q.Display()) < 10)
	q.Skip()
	q.Wait()
	assert.Equal(t, "東京1東京2東京3東京4東京5東京6東京7東京8\n東京9東京10東京11東京12", q.Display())
	q.Pop()
	q.Wait()
	assert.Equal(t, "last", q.Display())
}

func TestMsgEmit_流れる(t *testing.T) {
	t.Skip("未実装")
	input := `*start
あ
い
う
え
お`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program, err := p.ParseProgram()
	assert.NoError(t, err)
	e := NewEvaluator()
	e.Eval(program)
	q := NewQueue(e)
	q.Start()

	q.Run()
	q.Wait()
	assert.Equal(t, "あいうえお", q.Display())
}

func TestMsgEmit_イベントを消費する(t *testing.T) {
	t.Skip("未実装")
	evaluator := Evaluator{}
	q := NewQueue(&evaluator)
	q.Evaluator.Events = append(q.Evaluator.Events,
		utils.GetPtr(NewMsgEmit("あい")),
		utils.GetPtr(NewMsgEmit("うえ")),
		utils.GetPtr(NewMsgEmit("おか")),
	)
	q.Start()

	q.Run()
	q.Wait()
	assert.Equal(t, "あいうえおか", q.Display())
}

func TestRun_RunがPopとSkipを使い分ける(t *testing.T) {
	evaluator := Evaluator{}
	q := NewQueue(&evaluator)
	q.Evaluator.Events = append(q.Evaluator.Events, utils.GetPtr(NewMsgEmit("東京1東京2東京3東京4東京5東京6東京7東京8東京9東京10東京11東京12")))
	q.Evaluator.Events = append(q.Evaluator.Events, &Flush{})
	q.Evaluator.Events = append(q.Evaluator.Events, utils.GetPtr(NewMsgEmit("last")))
	q.Start()

	assert.Equal(t, "", q.Display())
	q.Run() // skip
	q.Wait()
	assert.Equal(t, "東京1東京2東京3東京4東京5東京6東京7東京8\n東京9東京10東京11東京12", q.Display())
	q.Run() // pop
	q.Wait()
	assert.Equal(t, "last", q.Display())
}

func TestRun_Skipを使わずに時間経過でも1回のRunで次のイベントに遷移する(t *testing.T) {
	evaluator := Evaluator{}
	q := NewQueue(&evaluator)
	q.Evaluator.Events = append(q.Evaluator.Events, utils.GetPtr(NewMsgEmit("あい")))
	q.Evaluator.Events = append(q.Evaluator.Events, utils.GetPtr(NewMsgEmit("うえ")))
	q.Start()

	time.Sleep(50 * time.Millisecond) // アニメーション時間経過
	assert.Equal(t, "あい", q.Display())
	q.Run() // pop
	q.Wait()
	assert.Equal(t, "あいうえ", q.Display())
}

func TestJump_複数実行できる(t *testing.T) {
	input := `*start
サンプル1
サンプル2[p]
新文章`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program, err := p.ParseProgram()
	assert.NoError(t, err)
	e := NewEvaluator()
	e.Eval(program)
	q := NewQueue(e)
	q.Start()

	assert.Equal(t, "", q.Display())
	q.Run() // run
	assert.Equal(t, "", q.Display())
	q.Wait()
	assert.Equal(t, "サンプル1", q.Display())
	q.Run() // pop
	assert.Equal(t, "サンプル1", q.Display())
	q.Wait()
	assert.Equal(t, "サンプル1サンプル2", q.Display())
	q.Run()
	q.Wait()
	assert.Equal(t, "新文章", q.Display())
}

func TestJump_ラベルジャンプできる(t *testing.T) {
	input := `*start
スタート[p]
[jump target="sample"]
*ignore
これは無視
*sample
サンプル1`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program, err := p.ParseProgram()
	assert.NoError(t, err)
	e := NewEvaluator()
	e.Eval(program)
	q := NewQueue(e)
	q.Start()

	assert.Equal(t, "", q.Display())
	q.Run() // skip
	assert.Equal(t, "", q.Display())
	q.Wait()
	assert.Equal(t, "スタート", q.Display())
	q.Run() // pop (->jump)
	assert.Equal(t, "スタート", q.Display())
	q.Wait()
	assert.Equal(t, "サンプル1", q.Display())
}

func TestWorker_startラベルから開始する(t *testing.T) {
	input := `*ignore
無視するべき
*start
スタート`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program, err := p.ParseProgram()
	assert.NoError(t, err)
	e := NewEvaluator()
	e.Eval(program)
	q := NewQueue(e)
	q.Start()

	q.Run()
	q.Wait()
	assert.Equal(t, "スタート", q.Display())
}

func TestImage_背景変更を通知する(t *testing.T) {
	input := `*start
[image source="test.png"]
スタート
[p]
ああああ
[p]`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program, err := p.ParseProgram()
	assert.NoError(t, err)
	e := NewEvaluator()
	e.Eval(program)
	q := NewQueue(e)
	q.Start()

	assert.Equal(t, "", q.Display())
	q.Run() // pop
	q.Wait()

	receivedEvent := <-q.NotifyChan
	assert.Equal(t, &ChangeBg{Source: "test.png"}, receivedEvent)

	assert.Equal(t, "スタート", q.Display())
	q.Run() // pop
	q.Wait()
	assert.Equal(t, "ああああ", q.Display())
}

// TODO: 一発で流れてほしい
func TestNewline_改行できる(t *testing.T) {
	input := `*start
あ[r]
い[r]
う[r]
え[r]
お[r]
[l]
かきくけこ
[p]`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program, err := p.ParseProgram()
	assert.NoError(t, err)
	e := NewEvaluator()
	e.Eval(program)
	q := NewQueue(e)
	q.Start()

	q.Run()
	q.Wait()
	q.Run()
	q.Wait()
	q.Run()
	q.Wait()
	q.Run()
	q.Wait()
	q.Run()
	q.Wait()
	assert.Equal(t, "あ\nい\nう\nえ\nお", q.Display())
}
