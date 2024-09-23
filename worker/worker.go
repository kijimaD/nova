package worker

import (
	"sync"
)

// queueて名前、おかしいかもしれない
// 文字列は構造体にしたい
type Queue struct {
	workerChan chan Event
	// イベント群
	Events []Event
	// 現在の表示文字列
	// アニメーション用に1文字ずつ増えていく
	buf string

	// 現在実行中
	cur Event

	wg sync.WaitGroup
}

func NewQueue() Queue {
	q := Queue{
		Events:     []Event{},
		workerChan: make(chan Event, 1024),
	}

	return q
}

// ワーカーを開始する
func (q *Queue) Start() {
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
	e := q.Events[0]
	q.cur = e
	q.wg.Add(1)
	q.workerChan <- e
	q.Events = append(q.Events[:0], q.Events[1:]...)

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
			if !ok {
				// closeしているので終了
				q.Pop()
			}
		default:
			q.Skip()
		}
	default:
		q.Pop()
		q.wg.Done()
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
	q.Events = es
}
