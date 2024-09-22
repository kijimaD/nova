package msg

import (
	"time"
)

type Event interface {
	Run(*Queue)
}

type Skipper interface {
	Skip()
}

// ================

type TaskStatus string

const (
	// 表示アニメーション中
	TaskRunning = TaskStatus("RUNNING")
	// メッセージをすべて表示した
	TaskFinish = TaskStatus("FINISH")
)

// メッセージ表示
type msgEmit struct {
	body []rune
	// 自動改行カウント
	nlCount int

	status    TaskStatus
	bChan     chan string
	IsAnimate bool
}

// 全部終わっていればPopする
// 終わってなければSkipする
// skipどうするか
func (e *msgEmit) Run(q *Queue) {
	e.status = TaskRunning
	e.bChan = make(chan string, 2048)
	e.IsAnimate = true

	go func() {
		for b := range e.bChan {
			q.buf += b
		}
	}()

	go func() {
		for {
			// skip := <-e.skipChan

			// finish := <-e.finishChan
		}
	}()

	// すでにRunが実行されていた場合には即時コピーしたい
	for {
		if len(e.body) == 0 {
			break
		}
		e.bChan <- string(e.body[0])
		e.body = e.body[1:]

		if e.IsAnimate {
			time.Sleep(10 * time.Millisecond)
		}
	}
	e.status = TaskFinish

	return
}

func (e *msgEmit) Skip() {
	e.IsAnimate = false
}

// ================

// ページをフラッシュする
type flush struct{}

func (c *flush) Run(q *Queue) {
	q.buf = ""
	return
}

// ================

type ChangeBg struct {
	Source string
}

func (c *ChangeBg) Run(q *Queue) {
	q.Pop()
	return
}

// ================

// 行末クリック待ち
type lineEndWait struct{}

func (l *lineEndWait) Run(q *Queue) {
	q.buf = q.buf + "\n"
	return
}

// ================

// 未実装
type notImplement struct{}

func (l *notImplement) Run(q *Queue) {
	q.buf = ""
	return
}

// ================
type wait struct {
	durationMsec time.Duration
}

func (w *wait) Run(q *Queue) {
	time.Sleep(w.durationMsec)
	q.buf = ""
	return
}
