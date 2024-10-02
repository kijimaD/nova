package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAutoNewline(t *testing.T) {
	assert.Equal(t,
		"",
		autoNewline("", 10),
	)
	assert.Equal(t,
		"あいうえ",
		autoNewline("あいうえ", 10),
	)
	assert.Equal(t,
		"あいうえおかきくけこ\nさしすせそ",
		autoNewline("あいうえおかきくけこさしすせそ", 10),
	)
	assert.Equal(t,
		"あいうえお\nかきくけこ\nさしすせそ",
		autoNewline("あいうえお\nかきくけこさしすせそ", 5),
	)
	assert.Equal(t,
		"あいうえお\nかきくけこ\nさしすせそ\nたちつてと",
		autoNewline("あいうえお\nかきくけこ\nさしすせそたちつてと", 5),
	)
	assert.Equal(t,
		"abcdefghij\nklmno",
		autoNewline("abcdefghijklmno", 10),
	)
	assert.Equal(t,
		"あいうえお\nかきくけこ\nさしすせそ",
		autoNewline("あいうえおかきくけこさしすせそ", 5),
	)
}
