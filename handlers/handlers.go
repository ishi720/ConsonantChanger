package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"myapp/module"
	"net/http"
	"net/url"

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

// VOICEVOX APIのベースURL
const voicevoxBaseURL = "http://localhost:50021"

// 音声生成リクエストの構造体
type VoiceRequest struct {
	Text string `json:"text"`
}

// 音声生成ハンドラ
func GenerateVoiceHandler(c echo.Context) error {
	// POSTリクエストからテキストを取得
	var req VoiceRequest
	if err := c.Bind(&req); err != nil {
		req.Text = c.QueryParam("text")
	}

	if req.Text == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "text parameter is required",
		})
	}

	// ずんだもんのspeaker ID (3)
	speakerID := 3

	// audio_query を取得
	params := url.Values{}
	params.Set("text", req.Text)
	params.Set("speaker", fmt.Sprintf("%d", speakerID))

	audioQueryURL := fmt.Sprintf("%s/audio_query?%s", voicevoxBaseURL, params.Encode())

	fmt.Printf("DEBUG: Requesting audio_query: %s\n", audioQueryURL)

	// POSTリクエストを送信
	resp, err := http.Post(audioQueryURL, "application/json", nil)
	if err != nil {
		fmt.Printf("ERROR: Failed to connect to VOICEVOX: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "VOICEVOXに接続できません。VOICEVOXが起動しているか確認してください。",
		})
	}
	defer resp.Body.Close()

	// レスポンスを読み込み
	audioQuery, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ERROR: Failed to read audio query response: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "音声クエリの読み取りに失敗しました",
		})
	}

	fmt.Printf("DEBUG: audio_query response status: %d, body length: %d\n", resp.StatusCode, len(audioQuery))

	// ステータスコードを確認
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("ERROR: VOICEVOX API error: %s\n", string(audioQuery))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("音声クエリの生成に失敗しました (status: %d): %s",
				resp.StatusCode, string(audioQuery)),
		})
	}

	// レスポンスが短すぎる場合はエラーとみなす
	if len(audioQuery) < 50 {
		fmt.Printf("ERROR: Response too short, likely an error: %s\n", string(audioQuery))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "VOICEVOXからのレスポンスが不正です: " + string(audioQuery),
		})
	}

	// audio_queryのJSONをパースしてデバッグ出力
	var audioQueryJSON map[string]interface{}
	if err := json.Unmarshal(audioQuery, &audioQueryJSON); err != nil {
		fmt.Printf("ERROR: Invalid JSON in audio_query response: %v\n", err)
		fmt.Printf("DEBUG: Response body: %s\n", string(audioQuery))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "VOICEVOXからの応答が不正です",
		})
	}

	fmt.Printf("DEBUG: audio_query successful, JSON keys: %v\n", getKeys(audioQueryJSON))

	// synthesis で音声を生成
	synthesisParams := url.Values{}
	synthesisParams.Set("speaker", fmt.Sprintf("%d", speakerID))
	synthesisURL := fmt.Sprintf("%s/synthesis?%s", voicevoxBaseURL, synthesisParams.Encode())

	fmt.Printf("DEBUG: Requesting synthesis: %s\n", synthesisURL)

	// POSTリクエストを送信
	resp2, err := http.Post(synthesisURL, "application/json", bytes.NewReader(audioQuery))
	if err != nil {
		fmt.Printf("ERROR: Failed to synthesize audio: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "音声合成に失敗しました: " + err.Error(),
		})
	}
	defer resp2.Body.Close()

	// ステータスコードを確認
	if resp2.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp2.Body)
		fmt.Printf("ERROR: Synthesis failed: %s\n", string(bodyBytes))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("音声合成に失敗しました (status: %d): %s",
				resp2.StatusCode, string(bodyBytes)),
		})
	}

	// 音声データを読み込み
	audioData, err := io.ReadAll(resp2.Body)
	if err != nil {
		fmt.Printf("ERROR: Failed to read audio data: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "音声データの読み取りに失敗しました",
		})
	}

	fmt.Printf("DEBUG: Audio data generated successfully, size: %d bytes\n", len(audioData))

	// WAVファイルとして返す
	c.Response().Header().Set("Content-Type", "audio/wav")
	c.Response().Header().Set("Content-Disposition", "inline; filename=zundamon.wav")
	return c.Blob(http.StatusOK, "audio/wav", audioData)
}

// ヘルパー関数: JSONのキーを取得
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
