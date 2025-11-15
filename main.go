package main

import (
	"html/template"
	"myapp/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echoインスタンスを作成
	e := echo.New()

	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 静的ファイル
	e.Static("/static", "static")

	// Template構造体初期化し、HTMLテンプレートをロード
	t := &handlers.Template{
		Templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = t

	// ルートハンドラを設定
	e.GET("/", handlers.RootHandler)

	// APIエンドポイントのルーティング
	e.GET("/api/getColockLanguage", handlers.GetColockLanguageHandler)

	// 音声生成APIエンドポイント
	e.GET("/api/generateVoice", handlers.GenerateVoiceHandler)
	e.POST("/api/generateVoice", handlers.GenerateVoiceHandler)

	// ポート8080でサーバーを起動
	e.Logger.Info("Starting server on :8080")
	e.Start(":8080")
}
