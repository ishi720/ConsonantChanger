package main

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"myapp/handlers"
)

func main() {
	// Echoインスタンスを作成
	e := echo.New()

	e.Static("/static", "static")

	// Template構造体初期化し、HTMLテンプレートをロード
	t := &handlers.Template{
		Templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = t

	// ルートハンドラを設定
	e.GET("/", handlers.RootHandler)

	// APIエンドポイントのルーティング（GETメソッド）
	e.GET("/api/getColockLanguage", handlers.GetColockLanguageHandler)

	// ポート8080でサーバーを起動
	e.Start(":8080")
}