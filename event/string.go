package event

import "strings"

// 直近の行を見て、横幅を超えていたら改行
func autoNewline(buf string, chunkSize int) string {
	split := strings.Split(buf, "\n")
	last := split[len(split)-1]

	var latestLine strings.Builder
	runes := []rune(last)
	for i, r := range runes {
		latestLine.WriteRune(r)
		// 文末の場合は改行を追加しない
		if (i+1)%chunkSize == 0 && i+1 != len(runes) {
			latestLine.WriteString("\n")
		}
	}

	var result string
	if len(split) > 1 {
		// 加工した末尾以外は元に戻す
		original := strings.Join(split[0:len(split)-1], "\n")
		result = original + "\n" + latestLine.String()
	} else {
		result = latestLine.String()
	}

	return result
}
