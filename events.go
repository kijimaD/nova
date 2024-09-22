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
	// 表示文字列に送信するバッファ
	bChan chan string
	// 終了
	doneChan chan bool
	// 自動改行カウント
	nlCount int
}

func (e *msgEmit) Run(q *Queue) {
	e.bChan = make(chan string, 2048)
	e.doneChan = make(chan bool, 1)

	go func() {
		for b := range e.bChan {
			q.buf += b
		}
	}()

	go func() {
		for i, char := range e.body {
			select {
			case <-e.doneChan:
				// フラグが立ったら残りの文字を一気に表示
				e.bChan <- e.body[i:]

				return
			default:
				// フラグが立ってないので1文字ずつ表示
				e.bChan <- string(char)
				time.Sleep(10 * time.Millisecond)
			}
		}
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
