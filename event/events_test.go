package event

import (
	"testing"
	"time"
	"unicode/utf8"

	"github.com/kijimaD/nova/utils"
	"github.com/stretchr/testify/assert"
)

func TestMsgEmit_Skipできる(t *testing.T) {
	q := prepareQueue(t, `*start
東京1東京2東京3東京4東京5東京6東京7東京8東京9東京10東京11東京12
[p]
last`)
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

func TestJump_複数実行できる(t *testing.T) {
	q := prepareQueue(t, `*start
サンプル1
サンプル2[p]
新文章`)
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
	q := prepareQueue(t, `*start
スタート[p]
[jump target="sample"]
*ignore
これは無視
*sample
サンプル1`)
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

func TestAutoNewline(t *testing.T) {
	assert.Equal(t,
		"",
		autoNewline("", 10),
	)
	assert.Equal(t,
		"あいうえ",
		autoNewline("あいうえ", 10),
	)
	assert.Equal(t,
		"あいうえおかきくけこ\nさしすせそ",
		autoNewline("あいうえおかきくけこさしすせそ", 10),
	)
	assert.Equal(t,
		"あいうえお\nかきくけこ\nさしすせそ",
		autoNewline("あいうえお\nかきくけこさしすせそ", 5),
	)
	assert.Equal(t,
		"あいうえお\nかきくけこ\nさしすせそ\nたちつてと",
		autoNewline("あいうえお\nかきくけこ\nさしすせそたちつてと", 5),
	)
	assert.Equal(t,
		"abcdefghij\nklmno",
		autoNewline("abcdefghijklmno", 10),
	)
	assert.Equal(t,
		"あいうえお\nかきくけこ\nさしすせそ",
		autoNewline("あいうえおかきくけこさしすせそ", 5),
	)
}

func TestImage_背景変更を通知する(t *testing.T) {
	q := prepareQueue(t, `*start
[image source="test.png"]
スタート
[p]
ああああ
[p]`)
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
	t.Skip()

	q := prepareQueue(t, `*start
あ[r]
い[r]
う[r]
え[r]
お[r]`)
	q.Start()

	q.Run()
	q.Wait()
	assert.Equal(t, "あ\nい\nう\nえ\nお\n", q.Display())
}
