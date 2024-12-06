package event

import (
	"math"
	"sync"
)

// queueて名前、おかしいかもしれない
// 文字列は構造体にしたい
type Queue struct {
	// 評価器
	Evaluator *Evaluator
	// 内部で利用するイベントキュー
	// すべてのイベントが入る可能性がある
	workerChan chan Event
	// クライアント側での実装が必要なイベントを通知するキュー
	// テキスト関係のイベントはbufに変換され、入らない
	NotifyChan chan Event
	// 現在表示中の文字列
	// 利用側はこの文字列を表示するだけで、いい感じに表示できる
	// アニメーション用に1文字ずつ増えていく
	buf string
	// 実行中イベント
	cur Event
	// WaitGroup
	wg sync.WaitGroup
	// アニメーション待ち状態かどうか
	OnAnim bool
}

func NewQueue(evaluator *Evaluator) Queue {
	q := Queue{
		Evaluator:  evaluator,
		workerChan: make(chan Event, 1024),
		NotifyChan: make(chan Event, 1024),
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

// 処理中インデックスを進める
func (q *Queue) Pop() Event {
	e := q.Events()[q.Evaluator.CurrentEventIdx]
	q.cur = e
	q.wg.Add(1)
	q.workerChan <- e
	q.Evaluator.CurrentEventIdx = int(math.Min(float64(len(q.Events())-1), float64(q.Evaluator.CurrentEventIdx+1)))

	return e
}

// 現在処理中のタスクをスキップする
func (q *Queue) Skip() {
	if e, ok := q.cur.(Skipper); ok {
		e.Skip()
	}
}

// デバッグ用
func (q *Queue) Reset() {
	q.Wait()
	q.buf = ""
	q.Evaluator.Play("start") // 各イベントのチャンネルがcloseしているので初期化する
	q.Pop()                   // 次イベントの先頭を読み込み

	return
}

// 実行中タスクに合わせてPop()もしくはSkip()する
// 入力待ちにならないイベント(画像表示とか)は、イベント実行時に自身でPop()するため、この分岐にはこない
func (q *Queue) Run() {
	q.OnAnim = false
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
