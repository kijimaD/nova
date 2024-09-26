package ast

import (
	"bytes"
	"strings"

	"github.com/kijimaD/nova/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// 構文解析器が生成する全てのASTのルートノードになる
type Program struct {
	Statements []Statement
}

// インターフェースで定義されている関数の1つ
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// インターフェースで定義されている関数の1つ
// 文字列表示してデバッグしやすいようにする
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// 式文は文に含まれる
type ExpressionStatement struct {
	Token      token.Token // 式の最初のトークン
	Expression Expression  // 式を保持
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type CmdExpression struct {
	Token      token.Token // '['トークン
	Expression Expression
}

func (ie *CmdExpression) expressionNode()      {}
func (ie *CmdExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *CmdExpression) String() string {
	var out bytes.Buffer

	out.WriteString("[")
	out.WriteString(ie.Expression.String())
	out.WriteString("]")

	return out.String()
}

// 本文
type TextLiteral struct {
	Token token.Token
	Value string
}

func (sl *TextLiteral) expressionNode()      {}
func (sl *TextLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *TextLiteral) String() string       { return sl.Token.Literal }

type CmdLiteral struct {
	Token      token.Token
	FuncName   Identifier
	Parameters NamedParams
}

func (fl *CmdLiteral) expressionNode()      {} // fnの結果をほかの変数に代入できたりするため。代入式の一部として扱うためには、式でないといけない
func (fl *CmdLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *CmdLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for k, v := range fl.Parameters.Map {
		params = append(params, k+"="+v)
	}

	out.WriteString("[")
	out.WriteString(fl.FuncName.Value)
	out.WriteString(strings.Join(params, ", "))
	out.WriteString("]")

	return out.String()
}

type Identifier struct {
	Token token.Token // token.IDENT トークン
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type NamedParams struct {
	Map map[string]string
}

func (n *NamedParams) expressionNode() {}
func (n *NamedParams) String() string {
	var out bytes.Buffer

	for k, v := range n.Map {
		out.WriteString(k)
		out.WriteString(" = ")
		out.WriteString(v)
	}
	return out.String()
}

// ラベルの本体となるブロック
// Statementと区別する必要はない気がするな
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// ラベル表現
type LabelLiteral struct {
	Token     token.Token // token.ASTERISKトークン
	LabelName Identifier
	Body      *BlockStatement // 定義の中身自体は文
}

func (le *LabelLiteral) expressionNode()      {}
func (le *LabelLiteral) TokenLiteral() string { return le.Token.Literal }
func (le *LabelLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(token.ASTERISK)
	out.WriteString(le.TokenLiteral())
	out.WriteString("\n")
	out.WriteString(le.Body.String())

	return out.String()
}
