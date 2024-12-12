package handlers

import (
	"html/template"
	"io"
	"myapp/module"
	"net/http"

	"github.com/labstack/echo/v4"
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
	return c.Render(http.StatusOK, "index.html", nil)
}

// APIエンドポイントのハンドラ
func GetColockLanguageHandler(c echo.Context) error {
	inputString := c.QueryParam("input_string")
	lineType := c.QueryParam("line_type")

	// 入力を処理して結果を取得
	result := module.ConsonantLockLanguage(inputString, lineType, false)

	// 結果をJSON形式で返す
	return c.JSON(http.StatusOK, Response{Result: result})
}
