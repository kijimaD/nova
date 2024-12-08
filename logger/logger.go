package logger

import (
	"log/slog"
	"os"
)

// TODO: ライブラリでロガーをどうすべきなのかよくわからない。ライブラリでなければmain.initで環境変数に応じて初期化、などするのだが
var MyLog *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
