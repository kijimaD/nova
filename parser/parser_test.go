package parser

import (
	"testing"

	"github.com/kijimaD/nova/ast"
	"github.com/kijimaD/nova/lexer"

	"github.com/stretchr/testify/assert"
)

func TestParseProgram(t *testing.T) {
	input := `こんにちは[l]世界[p]`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program, err := p.ParseProgram()
	assert.NoError(t, err)

	{
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)
		textLit, ok := stmt.Expression.(*ast.TextLiteral)
		assert.True(t, ok)
		assert.Equal(t, "こんにちは", textLit.Value)
	}
	{
		stmt, ok := program.Statements[1].(*ast.ExpressionStatement)
		assert.True(t, ok)
		cmdExp, ok := stmt.Expression.(*ast.CmdLiteral)
		assert.True(t, ok)
		assert.Equal(t, "[l]", cmdExp.String())
	}
	{
		stmt, ok := program.Statements[2].(*ast.ExpressionStatement)
		assert.True(t, ok)
		textLit, ok := stmt.Expression.(*ast.TextLiteral)
		assert.True(t, ok)
		assert.Equal(t, "世界", textLit.Value)
	}
	{
		stmt, ok := program.Statements[3].(*ast.ExpressionStatement)
		assert.True(t, ok)
		cmdExp, ok := stmt.Expression.(*ast.CmdLiteral)
		assert.True(t, ok)
		assert.Equal(t, "[p]", cmdExp.String())
	}

	assert.Equal(t, "こんにちは[l]世界[p]", program.String())
}

func TestParseCmdExpression(t *testing.T) {
	input := `[example a="value1" b="value2" c="test.png"]
hello, world[r]
hello, world
`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program, err := p.ParseProgram()
	assert.NoError(t, err)

	s := program.Statements[0]
	stmt, ok := s.(*ast.ExpressionStatement)
	assert.True(t, ok)

	f, ok := stmt.Expression.(*ast.CmdLiteral)
	assert.True(t, ok)
	assert.Equal(t, "example", f.FuncName.Value)
	assert.Equal(t, "value1", f.Parameters.Map["a"])
	assert.Equal(t, "value2", f.Parameters.Map["b"])
	assert.Equal(t, "test.png", f.Parameters.Map["c"])

	assert.Equal(t, "[example a=value1, b=value2, c=test.png]hello, world[r]hello, world", program.String())
}

func TestParseText(t *testing.T) {
	input := `hello, world1
hello, world2
[p]
hello, world3
`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program, err := p.ParseProgram()
	assert.NoError(t, err)

	assert.Equal(t, "hello, world1hello, world2[p]hello, world3", program.String())
}

func TestParseCmdExpression_シンタックスエラーを捕捉できる(t *testing.T) {
	{
		input := `[example a]`

		l := lexer.NewLexer(input)
		p := NewParser(l)
		_, err := p.ParseProgram()
		assert.Error(t, err)

		assert.Equal(t, 1, len(p.Errors()))
		assert.Contains(t, p.Errors()[0], "シンタックスエラー: EQUALがない")
	}
	{
		input := `[example a=]`

		l := lexer.NewLexer(input)
		p := NewParser(l)
		_, err := p.ParseProgram()
		assert.Error(t, err)

		assert.Equal(t, 1, len(p.Errors()))
		assert.Contains(t, p.Errors()[0], "STRINGがない")
	}
	{
		input := `[example a="hello"`

		l := lexer.NewLexer(input)
		p := NewParser(l)
		_, err := p.ParseProgram()
		assert.Error(t, err)

		assert.Equal(t, 1, len(p.Errors()))
		assert.Contains(t, p.Errors()[0], "対応する右ブラケットが存在しなかったため、末尾まで到達した")
	}
}

func TestParseLabelExpression(t *testing.T) {
	input := `*example1
本文1
本文2
*example2
本文3
本文4
*こんにちは
本文5`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program, err := p.ParseProgram()
	assert.NoError(t, err)

	{
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)
		f, ok := stmt.Expression.(*ast.LabelLiteral)
		assert.True(t, ok)
		assert.Equal(t, "example1", f.LabelName.Value)

		{
			es, ok := f.Body.Statements[0].(*ast.ExpressionStatement)
			assert.True(t, ok)
			assert.Equal(t, "本文1", es.Token.Literal)
		}
		{
			es, ok := f.Body.Statements[1].(*ast.ExpressionStatement)
			assert.True(t, ok)
			assert.Equal(t, "\n", es.Token.Literal)
		}
		{
			es, ok := f.Body.Statements[2].(*ast.ExpressionStatement)
			assert.True(t, ok)
			assert.Equal(t, "本文2", es.Token.Literal)
		}
	}
	{
		stmt, ok := program.Statements[1].(*ast.ExpressionStatement)
		assert.True(t, ok)
		f, ok := stmt.Expression.(*ast.LabelLiteral)
		assert.True(t, ok)
		assert.Equal(t, "example2", f.LabelName.Value)

		{
			es, ok := f.Body.Statements[0].(*ast.ExpressionStatement)
			assert.True(t, ok)
			assert.Equal(t, "本文3", es.Token.Literal)
		}
		{
			es, ok := f.Body.Statements[1].(*ast.ExpressionStatement)
			assert.True(t, ok)
			assert.Equal(t, "\n", es.Token.Literal)
		}
		{
			es, ok := f.Body.Statements[2].(*ast.ExpressionStatement)
			assert.True(t, ok)
			assert.Equal(t, "本文4", es.Token.Literal)
		}
	}
	{
		stmt, ok := program.Statements[2].(*ast.ExpressionStatement)
		assert.True(t, ok)
		f, ok := stmt.Expression.(*ast.LabelLiteral)
		assert.True(t, ok)
		assert.Equal(t, "こんにちは", f.LabelName.Value)

		{
			es, ok := f.Body.Statements[0].(*ast.ExpressionStatement)
			assert.True(t, ok)
			assert.Equal(t, "本文5", es.Token.Literal)
		}
	}
}
