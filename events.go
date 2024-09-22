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
	// パーサーから渡ってきた表示対象の文字列
	body string
	// 終了
	doneChan chan bool
	// 自動改行カウント
	nlCount int
}

func (e *msgEmit) Run(q *Queue) {
	e.doneChan = make(chan bool, 1)

	go func() {
		for i, char := range e.body {
			select {
			case <-e.doneChan:
				// フラグが立ったら残りの文字を一気に表示
				q.buf += e.body[i:]
				q.wg.Done()

				return
			default:
				// フラグが立ってないので1文字ずつ表示
				q.buf += string(char)
				time.Sleep(10 * time.Millisecond)
			}
		}
		q.wg.Done()
	}()

	return
}

func (e *msgEmit) Skip() {
	e.doneChan <- true
	close(e.doneChan)
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
