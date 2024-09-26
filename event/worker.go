package event

import (
	"log"
	"sync"
)

// queueて名前、おかしいかもしれない
// 文字列は構造体にしたい
type Queue struct {
	workerChan chan Event
	// 評価器
	Evaluator *Evaluator
	// 現在の表示文字列
	// アニメーション用に1文字ずつ増えていく
	buf string

	// 現在実行中
	cur Event

	wg sync.WaitGroup
}

func NewQueue(evaluator *Evaluator) Queue {
	q := Queue{
		Evaluator:  evaluator,
		workerChan: make(chan Event, 1024),
	}

	return q
}

func (q *Queue) Events() []Event {
	return q.Evaluator.Events
}

// ワーカーを開始する
func (q *Queue) Start() {
	q.Evaluator.Play("start") // startラベルで開始する
	go func() {
		for {
			select {
			case event := <-q.workerChan:
				event.Run(q)
			}
		}
	}()

	q.Pop()
}

// 未処理キューの先頭を取り出して処理キューに入れる
func (q *Queue) Pop() Event {
	e := q.Events()[0]
	q.cur = e
	q.wg.Add(1)
	q.workerChan <- e
	q.Evaluator.Events = append(q.Events()[:0], q.Events()[1:]...)

	return e
}

// 現在処理中のタスクをスキップする
func (q *Queue) Skip() {
	if e, ok := q.cur.(Skipper); ok {
		e.Skip()
	}
}

// 実行中タスクに合わせてPop()もしくはSkip()する
func (q *Queue) Run() {
	switch v := q.cur.(type) {
	case *MsgEmit:
		select {
		case _, ok := <-v.DoneChan:
			if ok {
				q.Pop()
			}
		default:
			// チャネルがクローズされているわけでもなく、値もまだ来ていない
			q.Skip()
		}
	default:
		log.Printf("想定してないイベントタイプ: %s", v)
	}
}

// すべてのジョブが処理されるまで待機
func (q *Queue) Wait() {
	q.wg.Wait()
}

// 処理中タスクを取得する
func (q *Queue) Head() Event {
	return q.cur
}

func (q *Queue) Display() string {
	return q.buf
}

func (q *Queue) SetEvents(es []Event) {
	q.Evaluator.Events = es
}
