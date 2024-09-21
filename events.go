package msg

import (
	"time"
)

type Event interface {
	Run(*Queue)
}

// ================

type TaskStatus string

const (
	// 未実行
	TaskNotRunning = TaskStatus("NOT_RUNNING")
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

	status TaskStatus
	bChan  chan string
}

// 全部終わっていればPopする
// 終わってなければSkipする
func (e *msgEmit) Run(q *Queue) {
	e.bChan = make(chan string)

	const width = 14

	go func() {
		for b := range e.bChan {
			q.buf += b
		}
	}()

	// 時間でbuf->allへのコピーを実行する
	// すでにRunが実行されていた場合には即時コピーする
	go func() {
		for {
			if len(e.body) == 0 {
				break
			}
			e.bChan <- string(e.body[0])
			e.body = e.body[1:]
			// 意図的に挿入された改行がある場合はカウンタをリセット
			// if string(e.body[e.pos]) == "\n" {
			// 	e.nlCount = 0
			// }
			// if e.nlCount%width == width-1 {
			// 	q.buf += "\n"
			// }
			// e.pos++
			// e.nlCount++

			if !(e.status == TaskFinish) {
				// time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return
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
	q.Pop()
	return
}

// ================

// 未実装
type notImplement struct{}

func (l *notImplement) Run(q *Queue) {
	q.buf = ""
	q.Pop()
	return
}

// ================
type wait struct {
	durationMsec time.Duration
}

func (w *wait) Run(q *Queue) {
	time.Sleep(w.durationMsec)
	q.buf = ""
	q.Pop()
	return
}
