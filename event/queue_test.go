package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPlay_指定ラベルを読み込める(t *testing.T) {
	q := prepareQueue(t, `*start
xxx`)
	err := q.Play("start")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(q.EventQueue))
}

func TestPlay_指定ラベルが存在しないとエラーを返す(t *testing.T) {
	q := prepareQueue(t, `*start
xxx`)
	err := q.Play("not exists")
	assert.Error(t, err)
	assert.Equal(t, 0, len(q.EventQueue))
}

func TestRun_RunがPopとSkipを使い分ける(t *testing.T) {
	q := prepareQueue(t, `*start
東京1東京2東京3東京4東京5東京6東京7東京8東京9東京10東京11東京12
[p]
last
`)
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
	q := prepareQueue(t, `*start
あい
うえ`)
	q.Start()

	time.Sleep(50 * time.Millisecond) // アニメーション時間経過
	assert.Equal(t, "あい", q.Display())
	q.Run() // pop
	q.Wait()
	assert.Equal(t, "あいうえ", q.Display())
}

func TestWorker_startラベルから開始する(t *testing.T) {
	q := prepareQueue(t, `*ignore
無視するべき
*start
スタート`)
	q.Start()

	q.Run()
	q.Wait()
	assert.Equal(t, "スタート", q.Display())
}
