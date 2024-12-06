package event

import (
	"log"
	"strings"
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

var (
	messageSpeed = 20 * time.Millisecond
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
		case _, ok := <-e.DoneChan:
			// フラグが立ったら残りの文字を一気に表示
			if ok {
				q.buf += e.Body[i:]
				q.buf = autoNewline(q.buf, 24)
				q.wg.Done()
			}
			// FIXME: チェックによってチャンネルの値を消費したが、workerのselect文で必要なので再度通知する...
			// closeしたほうがいいのかもしれないが、closeがかぶることがあり、その回避のためコードがわかりにくくなるので、再度通知を送ることにした
			e.DoneChan <- true
			q.OnAnim = true

			return
		default:
			// フラグが立ってないので1文字ずつ表示
			q.buf += string(char)
			q.buf = autoNewline(q.buf, 24)
			time.Sleep(messageSpeed)
		}
	}

	// 1文字ずつ表示し終わった場合
	e.DoneChan <- true
	q.OnAnim = true
	q.wg.Done()

	return
}

// 直近の行を見て、横幅を超えていたら改行
func autoNewline(buf string, chunkSize int) string {
	split := strings.Split(buf, "\n")
	last := split[len(split)-1]

	var latestLine strings.Builder
	runes := []rune(last)
	for i, r := range runes {
		latestLine.WriteRune(r)
		// 文末の場合は改行を追加しない
		if (i+1)%chunkSize == 0 && i+1 != len(runes) {
			latestLine.WriteString("\n")
		}
	}

	var result string
	if len(split) > 1 {
		// 加工した末尾以外は元に戻す
		original := strings.Join(split[0:len(split)-1], "\n")
		result = original + "\n" + latestLine.String()
	} else {
		result = latestLine.String()
	}

	return result
}

func (e *MsgEmit) Skip() {
	e.DoneChan <- true
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

// 行末クリック待ち
type LineEndWait struct{}

func (l *LineEndWait) Run(q *Queue) {
	q.buf += "\n"
	q.Pop()
	q.wg.Done()
	return
}

// ================

// 背景変更待ち
type ChangeBg struct {
	Source string
}

func (c *ChangeBg) Run(q *Queue) {
	q.Pop()
	q.wg.Done()
	q.NotifyChan <- c
	return
}

// ================

// 秒数待ち
type Wait struct {
	DurationMsec time.Duration
}

func (w *Wait) Run(q *Queue) {
	time.Sleep(w.DurationMsec)
	q.Pop()
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

type Newline struct{}


func (n *Newline) Run(q *Queue) {
	q.buf += "\n"
	q.Pop()
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
