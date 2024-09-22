package msg

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
}

func NewQueue(events []Event) Queue {
	q := Queue{
		events:     events,
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

	return NewQueue(e.Events)
}

func (q *Queue) Start() {
	go func() {
		for {
			select {
			case event := <-q.workerChan:
				q.cur = event
				event.Run(q)
			}
		}
	}()
}

// 未処理キューの先頭を取り出して処理キューに入れる
func (q *Queue) Pop() Event {
	e := q.events[0]
	q.workerChan <- e
	q.events = append(q.events[:0], q.events[1:]...)

	return e
}

// 現在処理中のタスクをスキップする
func (q *Queue) Skip() {
	if e, ok := q.cur.(Skipper); ok {
		e.Skip()
	}
}

// タスクに合わせてPop()もしくはSkip()する
func (q *Queue) Run() {
	switch v := q.cur.(type) {
	case *msgEmit:
		switch v.status {
		case TaskRunning:
			q.Skip()
		case TaskFinish:
			q.Pop()
		}
	default:
		q.Pop()
	}
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
