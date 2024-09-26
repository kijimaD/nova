package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAutoNewline(t *testing.T) {
	assert.Equal(t, "あいうえおかきくけこ\nさしすせそ", autoNewline("あいうえおかきくけこさしすせそ", 10))
	assert.Equal(t, "abcdefghij\nklmno", autoNewline("abcdefghijklmno", 10))
	assert.Equal(t, "あいうえお\nかきくけこ\nさしすせそ", autoNewline("あいうえおかきくけこさしすせそ", 5))
}
