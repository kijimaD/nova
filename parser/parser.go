package parser

import (
	"fmt"
	"log"

	"github.com/kijimaD/nov/token"

	"github.com/kijimaD/nov/ast"
	"github.com/kijimaD/nov/lexer"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token // 現在のトークン
	peekToken token.Token // 次のトークン

	// 構文解析関数
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	// 前置構文解析関数。前置演算子には「左側」が存在しない
	prefixParseFn func() ast.Expression
	// 中置構文解析関数 n + 1
	// 引数は中置演算子の「左側」
	infixParseFn func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	CMD // [...]
)

// 優先順位テーブル。トークンタイプと優先順位を関連付ける
var precedences = map[token.TokenType]int{
	token.LBRACKET: CMD,
}

// 字句解析器を受け取って初期化する
func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// 前置トークン
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.TEXT, p.parseTextLiteral)
	p.registerPrefix(token.LBRACKET, p.parseCmdLiteral)

	// 2つトークンを読み込む。curTokenとpeekTokenの両方がセットされる
	p.nextToken()
	p.nextToken()

	return p
}

// エラーのアクセサ
func (p *Parser) Errors() []string {
	return p.errors
}

// エラーを追加する
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t,
		p.peekToken.Type,
	)
	p.errors = append(p.errors, msg)
}

// 次のトークンに進む
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// パースを開始する。トークンを1つずつ辿る
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// 文をパースする。トークンの型によって適用関数を変える
func (p *Parser) parseStatement() ast.Statement {
	// 式文の構文解析を試みる
	return p.parseExpressionStatement()
}

// 式文を構文解析する
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

// 現在のトークンと引数の型を比較する
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// 次のトークンと引数の型を比較する
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// peekTokenの型をチェックし、その型が正しい場合に限ってnextTokenを読んで、トークンを進める
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// 構文解析関数を登録する
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	// 優先順位の処理
	// より低い優先順位のトークンに遭遇する間繰り返す
	// 優先順位が同じもしくは高いトークンに遭遇すると実行しない
	for precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

// 次のトークンタイプに対応している優先順位を返す
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// 文字列トークンをパース
func (p *Parser) parseTextLiteral() ast.Expression {
	return &ast.TextLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// コマンドリテラルをパース
// [image storage="test.png"]
// [p]
func (p *Parser) parseCmdLiteral() ast.Expression {
	lit := &ast.CmdLiteral{Token: p.curToken}
	p.nextToken()
	ident := ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	lit.FuncName = ident

	if !p.peekTokenIs(token.RBRACKET) {
		p.nextToken()
	}
	lit.Parameters = p.parseCmdParameters()

	p.nextToken()

	return lit
}

// 引数をパース
func (p *Parser) parseCmdParameters() ast.NamedParams {
	namedParams := ast.NamedParams{}
	namedParams.Map = map[string]string{}

	for !p.peekTokenIs(token.RBRACKET) {
		name := ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		if !p.peekTokenIs(token.EQUAL) {
			log.Fatal("シンタックスエラー: EQUALがない: ", p.curToken.Literal)
		}
		p.nextToken()
		if !p.peekTokenIs(token.STRING) {
			log.Fatal("シンタックスエラー: STRINGがない: ", p.curToken.Literal)
		}
		p.nextToken()
		namedParams.Map[name.Value] = p.curToken.Literal

		if p.peekTokenIs(token.RBRACKET) {
			break
		}
		if p.peekTokenIs(token.EOF) {
			log.Fatal("対応する右ブラケットが存在しないため、末尾まで到達した")
		}

		p.nextToken()
	}

	return namedParams
}
