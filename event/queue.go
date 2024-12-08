package event

import (
	"fmt"
	"log"
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
	// このチャンネルに送られると待ちキューの先頭からイベントを読み込み、待ちキューの先頭を削除する
	popChan chan struct{}
	// 現在表示中の文字列
	// 利用側はこの文字列を表示するだけで、いい感じに表示できる
	// アニメーション用に1文字ずつ増えていく
	buf string
	// 実行中イベント
	cur Event
	// WaitGroup。クリック待ちイベントに到達するまでで1つの単位としている。そのため複数グループをまたぐ可能性がある
	// テストでWait()して確認しやすくする
	wg sync.WaitGroup
	// アニメーション待ち状態かどうか
	OnAnim bool

	// 現在実行中のラベル。クライアントが再生中のラベルを表示するのに使う
	CurrentLabel string
	// 実行待ちのイベントキュー。ここにある時点ではまだ実行されているわけではない
	// TODO: 名前をそれとわかるものに変更する
	EventQueue []Event
}

func NewQueue(evaluator *Evaluator) Queue {
	q := Queue{
		Evaluator:  evaluator,
		workerChan: make(chan Event, 1024),
		NotifyChan: make(chan Event, 1024),
		popChan:    make(chan struct{}, 1),
	}

	return q
}

// 処理待受を開始する
func (q *Queue) Start() {
	err := q.Play("start") // startラベルで開始する
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			select {
			case event := <-q.workerChan:
				event.Before(q)

				_, isSkip := event.(Skipper)
				// スキップ可能イベントはevent内で処理するのでここでは処理がいらない
				if !isSkip {
					// クリック待ちするイベントではDoneを発行する
					_, isBlock := event.(Blocker)
					if isBlock {
						q.wg.Done()
					} else {
						q.popChan <- struct{}{}
						fmt.Println("popChan通知@notIsWait")
					}
				}
			}
		}
	}()

	go func() {
		for range q.popChan {
			q.Pop()
		}
	}()

	q.wg.Add(1)
	// 初回Popは初期値を確実にセットするために即時実行する
	q.Pop()
	fmt.Println("popChan通知@初回")
}

func (q *Queue) Play(label string) error {
	q.CurrentLabel = label
	err := q.Evaluator.Play(label)
	if err != nil {
		return err
	}

	newQueue := make([]Event, len(q.Evaluator.Events))
	copy(newQueue, q.Evaluator.Events)
	q.EventQueue = newQueue

	return nil
}

// 処理中インデックスを進める
// イベント列の先頭をチャンネルに入れて、現在処理中とする。そして処理したイベント列の先頭を切る
// 名前から想像する挙動は、切り出してからイベントに入れる、であるが...
func (q *Queue) Pop() {
	if len(q.EventQueue) > 0 {
		q.cur = q.EventQueue[0]
		q.workerChan <- q.cur
		q.EventQueue = q.EventQueue[1:]
	}
}

// 現在処理中の、スキップ可能なタスクをスキップする
func (q *Queue) Skip() {
	if e, ok := q.cur.(Skipper); ok {
		e.Skip()
	}
}

// クリックを押したときに実行される想定
// 実行中タスクに合わせてPop()もしくはSkip()する
// 非ブロックのイベントでは、自動でPopするのでこの関数を通過しない
func (q *Queue) Run() {
	q.OnAnim = false
	q.cur.After(q)
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

// for debug
func (q *Queue) DumpQueue() []string {
	result := []string{}
	for _, e := range q.EventQueue {
		result = append(result, e.String())
	}

	return result
}
