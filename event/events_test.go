package event

import (
	"testing"
	"time"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestMsgEmit_Skipできる(t *testing.T) {
	q := prepareQueue(t, `*start
first[l]
[p]
last[l]
`)
	q.Start()

	assert.Equal(t, "", q.Display())
	time.Sleep(30 * time.Millisecond)
	assert.True(t, utf8.RuneCountInString(q.Display()) > 1)
	assert.True(t, utf8.RuneCountInString(q.Display()) < 10)
	q.Skip()
	q.Wait()
	assert.Equal(t, "first", q.Display())
	q.Run()
	q.Wait()
	assert.Equal(t, "last", q.Display())
}

func TestMsgEmit_Skipを使わずに時間経過でも表示できる(t *testing.T) {
	q := prepareQueue(t, `*start
あい
[l]
うえ
[l]`)
	q.Start()

	time.Sleep(50 * time.Millisecond) // アニメーション時間経過
	assert.Equal(t, "あい", q.Display())
	q.Run()
	time.Sleep(50 * time.Millisecond) // アニメーション時間経過
	assert.Equal(t, "あいうえ", q.Display())
}

func TestJump_ラベルジャンプできる(t *testing.T) {
	q := prepareQueue(t, `*start
スタート[l][p]
[jump target="sample"]
*ignore
これは無視
*sample
サンプル1[l]`)
	q.Start()

	assert.Equal(t, "", q.Display())
	q.Run()
	assert.Equal(t, "", q.Display())
	q.Wait()
	assert.Equal(t, "スタート", q.Display())
	q.Run()
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
[l]
[p]
ああああ
[l]
[p]`)
	q.Start()

	assert.Equal(t, "", q.Display())
	q.Run()
	q.Wait()

	receivedEvent := <-q.NotifyChan
	assert.Equal(t, &ChangeBg{Source: "test.png"}, receivedEvent)

	assert.Equal(t, "スタート", q.Display())
	q.Run()
	q.Wait()
	assert.Equal(t, "ああああ", q.Display())
}

func TestNewline_改行できる(t *testing.T) {
	q := prepareQueue(t, `*start
あああ[r][r][r]ううう[l][p][r]えええ[r]おおお[r][l]`)
	q.Start()

	q.Run()
	q.Wait()
	assert.Equal(t, "あああ\n\n\nううう", q.Display())
	q.Run()
	q.Wait()
	assert.Equal(t, "\nえええ\nおおお\n", q.Display())
}
