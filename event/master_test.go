package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLabel_keyが存在すると返す(t *testing.T) {
	q := prepareQueue(t, `*start
xxx`)
	lm := q.Evaluator.LabelMaster
	label, err := lm.GetLabel("start")
	assert.NoError(t, err)

	assert.Equal(t, "start", label.Name)
	body := *label.Body
	assert.Equal(t, "xxx", body.String())
}

func TestGetLabel_keyが存在しないとエラーを返す(t *testing.T) {
	q := prepareQueue(t, `*start
xxx`)
	lm := q.Evaluator.LabelMaster
	_, err := lm.GetLabel("not exists")
	assert.Error(t, err)
}
