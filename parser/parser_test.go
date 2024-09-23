package parser

import (
	"testing"

	"github.com/kijimaD/nov/ast"
	"github.com/kijimaD/nov/lexer"

	"github.com/stretchr/testify/assert"
)

func TestParsingIndexExpressions(t *testing.T) {
	input := `こんにちは[l]世界[p]
←無視される改行たたたたた。
←有効な改行`

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
		cmdExp, ok := stmt.Expression.(*ast.FunctionLiteral)
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
		cmdExp, ok := stmt.Expression.(*ast.FunctionLiteral)
		assert.True(t, ok)
		assert.Equal(t, "[p]", cmdExp.String())
	}
	{
		stmt, ok := program.Statements[4].(*ast.ExpressionStatement)
		assert.True(t, ok)
		textLit, ok := stmt.Expression.(*ast.TextLiteral)
		assert.True(t, ok)
		assert.Equal(t, "←無視される改行たたたたた。\n←有効な改行", textLit.Value)
	}
}

func TestParsingCmdExpressionImage(t *testing.T) {
	input := `[image a="value1" b="value2" c="test.png"]`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	s := program.Statements[0]
	stmt, ok := s.(*ast.ExpressionStatement)
	assert.True(t, ok)

	f, ok := stmt.Expression.(*ast.FunctionLiteral)
	assert.True(t, ok)
	assert.Equal(t, "image", f.FuncName.Value)
	assert.Equal(t, "value1", f.Parameters.Map["a"])
	assert.Equal(t, "value2", f.Parameters.Map["b"])
	assert.Equal(t, "test.png", f.Parameters.Map["c"])
}
