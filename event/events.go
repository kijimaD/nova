package event

import (
	"fmt"
	"strings"
	"time"
)

type Event interface {
	Run(*Queue)
	String() string
}

// アニメーション状態を持ち、スキップ可能なイベント
type Skipper interface {
	Skip()
}

// ================

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

func (e *MsgEmit) String() string {
	return fmt.Sprintf("<MsgEmit %s>", e.Body)
}

// 文字送り中か文字表示完了かの2通りの状態がある
// 表示終了したら完了チャンネルにフラグを送る
func (e *MsgEmit) Run(q *Queue) {
	lineLen := 24

	for i, char := range e.Body {
		select {
		case _, ok := <-e.DoneChan:
			// フラグが立ったら残りの文字を一気に表示
			if !ok {
				continue
			}
			q.buf += e.Body[i:]
			q.buf = autoNewline(q.buf, lineLen)
			q.wg.Done()
			close(e.DoneChan)
			q.OnAnim = true

			return
		default:
			// フラグが立ってないので1文字ずつ表示
			q.buf += string(char)
			q.buf = autoNewline(q.buf, lineLen)
			time.Sleep(messageSpeed)
		}
	}

	// 1文字ずつ表示し終わった場合
	e.DoneChan <- true
	q.OnAnim = true
	q.wg.Done()

	q.Pop()

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

func (c *Flush) String() string {
	return "<Flush>"
}

func (c *Flush) Run(q *Queue) {
	q.buf = ""

	return
}

// ================

// 行末クリック待ち
type LineEndWait struct{}

func (l *LineEndWait) String() string {
	return "<LineEndWait>"
}

func (l *LineEndWait) Run(q *Queue) {
	q.buf += "\n"

	return
}

// ================

// 背景変更待ち
type ChangeBg struct {
	Source string
}

func (c *ChangeBg) String() string {
	return fmt.Sprintf("<ChangeBg %s>", c.Source)
}

func (c *ChangeBg) Run(q *Queue) {
	q.NotifyChan <- c

	return
}

// ================

// 秒数待ち
type Wait struct {
	DurationMsec time.Duration
}

func (w *Wait) String() string {
	return fmt.Sprintf("<Wait %s>", w.DurationMsec)
}

func (w *Wait) Run(q *Queue) {
	time.Sleep(w.DurationMsec)

	return
}

// ================

// ジャンプ。別のラベルへ遷移する
type Jump struct {
	Target string
}

func (j *Jump) String() string {
	return fmt.Sprintf("<Jump %s>", j.Target)
}

func (j *Jump) Run(q *Queue) {
	q.Play(j.Target)

	return
}

type Newline struct{}

func (n *Newline) String() string {
	return "<Newline>"
}

func (n *Newline) Run(q *Queue) {
	q.buf += "\n"

	return
}

// ================

// 未実装
type NotImplement struct{}

func (l *NotImplement) String() string {
	return "NotImplement"
}

func (l *NotImplement) Run(q *Queue) {
	return
}
