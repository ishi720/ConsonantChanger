package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	// Echoインスタンスを作成
	e := echo.New()

	// ルートハンドラを設定
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Echo!")
	})

	// サーバーをポート8080で起動
	e.Start(":8080")
}