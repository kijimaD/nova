package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlay_指定ラベルを読み込める(t *testing.T) {
	q := prepareQueue(t, `*start
xxx
[p]
*example1
yyy
[p]
*example2
zzz
[p]`)
	q.Start()

	q.Run()
	q.Wait()
	assert.Equal(t, "xxx", q.Display())

	assert.NoError(t, q.Play("example1"))
	q.Run()
	q.Wait()
	assert.Equal(t, "yyy", q.Display())

	assert.NoError(t, q.Play("example2"))
	q.Run()
	q.Wait()
	assert.Equal(t, "zzz", q.Display())
}

func TestPlay_指定ラベルが存在しないとエラーを返す(t *testing.T) {
	q := prepareQueue(t, `*start
xxx`)
	err := q.Play("not exists")
	assert.Error(t, err)
	assert.Equal(t, 0, len(q.WaitingQueue))
}

func TestRun_PopとSkipを使い分ける(t *testing.T) {
	q := prepareQueue(t, `*start
xxx
[p]
last
[p]
`)
	q.Start()

	assert.Equal(t, "", q.Display())
	q.Run()
	q.Wait()
	assert.Equal(t, "xxx", q.Display())
	q.Run()
	q.Wait()
	assert.Equal(t, "last", q.Display())
}

func TestWorker_startラベルから開始する(t *testing.T) {
	q := prepareQueue(t, `*ignore
無視するべき[l]
*start
スタート[l]`)
	q.Start()

	q.Run()
	q.Wait()
	assert.Equal(t, "スタート", q.Display())
}
