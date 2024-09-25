package event

import (
	"log"
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
type MsgEmit struct {
	// パーサーから渡ってきた表示対象の文字列
	Body string
	// 終了
	DoneChan chan bool
}

func NewMsgEmit(body string) MsgEmit {
	return MsgEmit{
		Body:     body,
		DoneChan: make(chan bool, 1),
	}
}

func (e *MsgEmit) Run(q *Queue) {
	if e.DoneChan == nil {
		log.Fatal("doneChan is nil")
	}

	for i, char := range e.Body {
		select {
		case <-e.DoneChan:
			// フラグが立ったら残りの文字を一気に表示
			q.buf += e.Body[i:]
			q.wg.Done()

			return
		default:
			// フラグが立ってないので1文字ずつ表示
			q.buf += string(char)
			time.Sleep(20 * time.Millisecond)
		}
	}
	q.wg.Done()

	return
}

func (e *MsgEmit) Skip() {
	e.DoneChan <- true
	close(e.DoneChan)
}

// ================

// ページをフラッシュする
type Flush struct{}

func (c *Flush) Run(q *Queue) {
	q.buf = ""
	q.wg.Done()
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
type LineEndWait struct{}

func (l *LineEndWait) Run(q *Queue) {
	q.buf = q.buf + "\n"
	return
}

// ================

// 秒数待ち
type Wait struct {
	DurationMsec time.Duration
}

func (w *Wait) Run(q *Queue) {
	time.Sleep(w.DurationMsec)
	q.buf = ""
	return
}

// ================

// ジャンプ。別のラベルへ遷移する
type Jump struct {
	Target string
}

func (j *Jump) Run(q *Queue) {
	return
}

// ================

// 未実装
type NotImplement struct{}

func (l *NotImplement) Run(q *Queue) {
	q.buf = ""
	return
}
