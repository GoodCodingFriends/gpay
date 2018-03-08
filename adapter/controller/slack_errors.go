package controller

var errToSlackMessage = map[error]string{
	ErrUnknownCommand: "わからないコマンドだよ",
	ErrInvalidUsage:   "多分変な使い方してる",
	ErrInvalidUserID:  "@ をつけてもう一度お願い",
}
