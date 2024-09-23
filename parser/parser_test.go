package parser

import (
	"testing"

	"github.com/kijimaD/nov/ast"
	"github.com/kijimaD/nov/lexer"

	"github.com/stretchr/testify/assert"
)

func TestParsingIndexExpressions(t *testing.T) {
	input := `こんにちは[l]世界[p]`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

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
}

func TestParsingCmdExpression(t *testing.T) {
	input := `[example a="value1" b="value2" c="test.png"]`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	s := program.Statements[0]
	stmt, ok := s.(*ast.ExpressionStatement)
	assert.True(t, ok)

	f, ok := stmt.Expression.(*ast.CmdLiteral)
	assert.True(t, ok)
	assert.Equal(t, "example", f.FuncName.Value)
	assert.Equal(t, "value1", f.Parameters.Map["a"])
	assert.Equal(t, "value2", f.Parameters.Map["b"])
	assert.Equal(t, "test.png", f.Parameters.Map["c"])
}

func TestParsingLabelExpression(t *testing.T) {
	input := `*example1
本文1
本文2`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

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
}
