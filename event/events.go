package event

import (
	"fmt"
	"time"

	"github.com/kijimaD/nova/logger"
)

type Event interface {
	// クリック前に実行するフック
	Before(*Queue)
	// クリック後に実行するフック
	After(*Queue)
	// デバッグ時に表示する文字列
	String() string
}

// アニメーション状態を持ち、スキップ可能なイベント
type Skipper interface {
	Skip()
}

// クリック待ちにするイベント
type Blocker interface {
	IsBlock()
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

// 実行中タスクに合わせてPop()もしくはSkip()する
// 文字送り中か文字表示完了かの2通りの状態がある
func (e *MsgEmit) Before(q *Queue) {
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

			close(e.DoneChan)
			q.OnAnim = true

			q.popChan <- struct{}{}
			logger.MyLog.Debug("popChan通知@スキップ")

			return
		default:
			// フラグが立ってないので1文字ずつ表示
			q.buf += string(char)
			q.buf = autoNewline(q.buf, lineLen)
			time.Sleep(messageSpeed)
		}
	}

	// 1文字ずつ表示し終わった場合
	close(e.DoneChan)
	q.OnAnim = true

	q.popChan <- struct{}{}
	logger.MyLog.Debug("popChan通知@順当")

	return
}

func (e *MsgEmit) After(q *Queue) {
	select {
	case _, ok := <-e.DoneChan:
		// close
		if !ok {
			q.popChan <- struct{}{}
			logger.MyLog.Debug("popChan通知@Run/MsgEmit")
		}
	default:
		// チャネルがクローズされているわけでもなく、値もまだ来ていない
		e.Skip()
	}
}

func (e *MsgEmit) Skip() {
	e.DoneChan <- true
}

// ================

// クリック待ちにして、クリックしたあとにフラッシュする
type Flush struct{}

func (c *Flush) String() string {
	return "<Flush>"
}

func (c *Flush) Before(q *Queue) {}

func (c *Flush) After(q *Queue) {
	q.buf = ""

	q.popChan <- struct{}{}
	logger.MyLog.Debug("popChan通知@Flush")
	q.wg.Add(1)
}

func (c *Flush) IsBlock() {}

// ================

// クリック待ちにして、クリックしたあとに改行する
type LineEndWait struct{}

func (l *LineEndWait) String() string {
	return "<LineEndWait>"
}

func (l *LineEndWait) Before(q *Queue) {}

func (l *LineEndWait) After(q *Queue) {
	q.buf += "\n"

	q.popChan <- struct{}{}
	logger.MyLog.Debug("popChan通知@LineEndWait")
	q.wg.Add(1)
}

func (l *LineEndWait) IsBlock() {}

// ================

// 背景変更
type ChangeBg struct {
	Source string
}

func (c *ChangeBg) String() string {
	return fmt.Sprintf("<ChangeBg %s>", c.Source)
}

func (c *ChangeBg) Before(q *Queue) {
	q.NotifyChan <- c

	return
}

func (c *ChangeBg) After(q *Queue) {}

// ================

// 秒数待ち
type Wait struct {
	DurationMsec time.Duration
}

func (w *Wait) String() string {
	return fmt.Sprintf("<Wait %s>", w.DurationMsec)
}

func (w *Wait) Before(q *Queue) {
	time.Sleep(w.DurationMsec)

	return
}

func (w *Wait) After(q *Queue) {}

// ================

// ジャンプ。別のラベルへ遷移する
type Jump struct {
	Target string
}

func (j *Jump) String() string {
	return fmt.Sprintf("<Jump %s>", j.Target)
}

func (j *Jump) Before(q *Queue) {
	q.Play(j.Target)

	return
}

func (j *Jump) After(q *Queue) {
	q.Play(j.Target)

	return
}

// ================

type Newline struct{}

func (n *Newline) String() string {
	return "<Newline>"
}

func (n *Newline) Before(q *Queue) {
	q.buf += "\n"

	return
}

func (n *Newline) After(q *Queue) {}

// ================

// 未実装
type NotImplement struct{}

func (l *NotImplement) String() string {
	return "NotImplement"
}

func (l *NotImplement) Before(q *Queue) {}

func (l *NotImplement) After(q *Queue) {}
