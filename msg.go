package msg

// queueて名前、おかしいかもしれない
// 文字列は構造体にしたい
type Queue struct {
	// イベント群
	events []Event
	// 現在の表示文字列
	// アニメーション用に1文字ずつ増えていく
	buf string
}

func NewQueue(events []Event) Queue {
	q := Queue{
		events: events,
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

// キューの先端を取り出して実行する
// popしないこともあるので、名前に合っていない
func (q *Queue) Pop() Event {
	e := q.events[0]

	switch v := e.(type) {
	case *msgEmit:
		switch v.status {
		case TaskNotRunning:
			// 未実行の場合実行する
			v.Run(q)
			v.status = TaskRunning
		case TaskRunning:
			v.status = TaskFinish
		case TaskFinish:
			// 終了
			q.events = append(q.events[:0], q.events[1:]...)
		}
	default:
		e.Run(q)
		// TODO: この書き方変な気がする
		// ここで先端を取り出すが、そうするとHeadは次を表すことになるので、おかしい
		q.events = append(q.events[:0], q.events[1:]...)
	}

	return e
}

// キューの先頭を表示だけする
// Head,いらないか
func (q *Queue) Head() Event {
	if len(q.events) == 0 {
		return &notImplement{}
	}

	return q.events[0]
}

func (q *Queue) Display() string {
	return q.buf
}

func (q *Queue) SetEvents(es []Event) {
	q.events = es
}
