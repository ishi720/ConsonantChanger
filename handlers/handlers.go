package handlers

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
	Templates *template.Template // フィールド名を大文字に変更
}

// テンプレートをレンダリング
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

// ルートハンドラ
func RootHandler(c echo.Context) error {
	data := map[string]interface{}{
		"Message": module.ConsonantLockLanguage("こんにちは", "pa"),
	}
	return c.Render(http.StatusOK, "index.html", data)
}

// APIエンドポイントのハンドラ
func GetColockLanguageHandler(c echo.Context) error {
	inputString := c.QueryParam("input_string")
	lineType := c.QueryParam("line_type")

	// 入力を処理して結果を取得
	result := module.ConsonantLockLanguage(inputString, lineType)

	// 結果をJSON形式で返す
	return c.JSON(http.StatusOK, Response{Result: result})
}