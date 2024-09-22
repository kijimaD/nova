package msg

import (
	"sync"
)

// queueて名前、おかしいかもしれない
// 文字列は構造体にしたい
type Queue struct {
	workerChan chan Event
	// イベント群
	events []Event
	// 現在の表示文字列
	// アニメーション用に1文字ずつ増えていく
	buf string

	// 現在実行中
	cur Event

	wg sync.WaitGroup
}

func NewQueue() Queue {
	q := Queue{
		events:     []Event{},
		workerChan: make(chan Event, 1),
	}
	return q
}

// スクリプトからキューを初期化する
func NewQueueFromText(text string) Queue {
	l := NewLexer(text)
	p := NewParser(l)
	program := p.ParseProgram()
	e := NewEvaluator(program)
	q := NewQueue()
	q.events = e.Events

	return q
}

func (q *Queue) Start() {
	go func() {
		for {
			select {
			case event := <-q.workerChan:
				event.Run(q)
			}
		}
	}()
}

// 未処理キューの先頭を取り出して処理キューに入れる
func (q *Queue) Pop() Event {
	e := q.events[0]
	q.cur = e
	q.workerChan <- e
	q.wg.Add(1)
	q.events = append(q.events[:0], q.events[1:]...)

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
	case *msgEmit:
		select {
		case _, ok := <-v.doneChan:
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
	q.events = es
}
