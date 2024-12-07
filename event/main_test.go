package event

import (
	"testing"

	"github.com/kijimaD/nova/lexer"
	"github.com/kijimaD/nova/parser"
	"github.com/stretchr/testify/assert"
)

func prepareQueue(t *testing.T, input string) Queue {
	t.Helper()

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program, err := p.ParseProgram()
	assert.NoError(t, err)
	e := NewEvaluator()
	e.Eval(program)
	q := NewQueue(e)

	return q
}
