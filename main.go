package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"myapp/module"
)

// テンプレートの構造体
type Template struct {
	templates *template.Template
}

// テンプレートをレンダリング
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// Echoインスタンスを作成
	e := echo.New()

	// Template構造体初期化し、HTMLテンプレートをロード
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = t

	// ルートハンドラを設定
	e.GET("/", func(c echo.Context) error {
		data := map[string]interface{}{
			"Message": module.ConsonantLockLanguage("こんにちは"),
		}
		return c.Render(http.StatusOK, "index.html", data)
	})

	// サーバーをポート8080で起動
	e.Start(":8080")
}