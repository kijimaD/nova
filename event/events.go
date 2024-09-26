package event

import (
	"log"
	"time"
)

// 別packageに移したいが、ここで参照があるためできない
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
	// 終了判定チャンネル。closeしてれば終了
	DoneChan chan bool
}

func NewMsgEmit(body string) MsgEmit {
	return MsgEmit{
		Body:     body,
		DoneChan: make(chan bool, 1),
	}
}

func (e *MsgEmit) Run(q *Queue) {
	// 初期化漏れ対策
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
	close(e.DoneChan)
	q.wg.Done()

	return
}

func (e *MsgEmit) Skip() {
	close(e.DoneChan)
}

// ================

// ページをフラッシュする
type Flush struct{}

func (c *Flush) Run(q *Queue) {
	q.buf = ""
	q.Pop()
	q.wg.Done()
	return
}

// ================

type ChangeBg struct {
	Source string
}

func (c *ChangeBg) Run(q *Queue) {
	q.Pop()
	q.wg.Done()
	return
}

// ================

// 行末クリック待ち
type LineEndWait struct{}

func (l *LineEndWait) Run(q *Queue) {
	q.buf = q.buf + "\n"
	q.wg.Done()
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
	q.wg.Done()
	return
}

// ================

// ジャンプ。別のラベルへ遷移する
type Jump struct {
	Target string
}

func (j *Jump) Run(q *Queue) {
	q.Evaluator.Play(j.Target)
	q.Pop() // 次イベントの先頭を読み込み
	q.wg.Done()
	return
}

// ================

// 未実装
type NotImplement struct{}

func (l *NotImplement) Run(q *Queue) {
	q.wg.Done()
	return
}
