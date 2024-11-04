package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"myapp/module"
)

// APIレスポンスの構造体
type Response struct {
	Result string `json:"result"`
}

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
			"Message": module.ConsonantLockLanguage("こんにちは", "pa"),
		}
		return c.Render(http.StatusOK, "index.html", data)
	})

	// APIエンドポイントのルーティング（GETメソッド）
	e.GET("/api/getColockLanguage", func(c echo.Context) error {
		inputString := c.QueryParam("input_string")
		lineType := c.QueryParam("line_type")

		// 入力を処理して結果を取得
		result := module.ConsonantLockLanguage(inputString, lineType)

		// 結果をJSON形式で返す
		return c.JSON(http.StatusOK, Response{Result: result})
	})

	// ポート8080でサーバーを起動
	e.Start(":8080")
}