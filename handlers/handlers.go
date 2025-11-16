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
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

// APIレスポンスの構造体
type Response struct {
	Result string `json:"result"`
}

// テンプレートの構造体
type Template struct {
	Templates *template.Template
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

// VOICEVOX APIのベースURLを環境変数から取得
func getVoicevoxBaseURL() string {
	url := os.Getenv("VOICEVOX_URL")
	if url == "" {
		// デフォルトはlocalhost
		url = "http://localhost:50021"
	}
	return url
}

// 音声生成リクエストの構造体
type VoiceRequest struct {
	Text string `json:"text"`
}

// 音声生成ハンドラ
func GenerateVoiceHandler(c echo.Context) error {
	startTime := time.Now()
	voicevoxBaseURL := getVoicevoxBaseURL()

	fmt.Printf("\n========================================\n")
	fmt.Printf("音声生成リクエスト開始: %s\n", startTime.Format("15:04:05"))
	fmt.Printf("VOICEVOX URL: %s\n", voicevoxBaseURL)
	fmt.Printf("========================================\n")

	// POSTリクエストからテキストを取得
	var req VoiceRequest
	if err := c.Bind(&req); err != nil {
		fmt.Printf("DEBUG: POST bodyのbindに失敗、GETパラメータを試行\n")
		req.Text = c.QueryParam("text")
	}

	fmt.Printf("受信したテキスト: \"%s\" (長さ: %d文字)\n", req.Text, len(req.Text))

	if req.Text == "" {
		fmt.Printf("ERROR: テキストが空です\n")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "text parameter is required",
		})
	}

	// ずんだもんのspeaker ID (3)
	speakerID := 3
	fmt.Printf("使用する話者ID: %d (ずんだもん)\n", speakerID)

	// Step 1: audio_query を取得
	fmt.Printf("\n[Step 1] audio_query の取得\n")
	params := url.Values{}
	params.Set("text", req.Text)
	params.Set("speaker", fmt.Sprintf("%d", speakerID))

	audioQueryURL := fmt.Sprintf("%s/audio_query?%s", voicevoxBaseURL, params.Encode())
	fmt.Printf("リクエストURL: %s\n", audioQueryURL)

	// VOICEVOXの接続テスト
	fmt.Printf("VOICEVOXへの接続テスト中...\n")
	testResp, testErr := http.Get(fmt.Sprintf("%s/speakers", voicevoxBaseURL))
	if testErr != nil {
		fmt.Printf("ERROR: VOICEVOXに接続できません: %v\n", testErr)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "VOICEVOXに接続できません。VOICEVOXが起動しているか確認してください。",
		})
	}
	testResp.Body.Close()
	fmt.Printf("✓ VOICEVOX接続OK\n")

	// audio_queryリクエスト送信
	fmt.Printf("audio_queryリクエスト送信中...\n")
	queryStartTime := time.Now()
	resp, err := http.Post(audioQueryURL, "application/json", nil)
	queryDuration := time.Since(queryStartTime)

	if err != nil {
		fmt.Printf("ERROR: audio_queryリクエスト失敗 (%v): %v\n", queryDuration, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "VOICEVOXへのリクエストに失敗しました: " + err.Error(),
		})
	}
	defer resp.Body.Close()

	fmt.Printf("✓ audio_queryレスポンス受信 (%v)\n", queryDuration)
	fmt.Printf("  ステータスコード: %d\n", resp.StatusCode)

	// レスポンスを読み込み
	audioQuery, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ERROR: レスポンス読み取り失敗: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "音声クエリの読み取りに失敗しました",
		})
	}

	fmt.Printf("  レスポンスサイズ: %d バイト\n", len(audioQuery))

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("ERROR: audio_query失敗 (status: %d)\n", resp.StatusCode)
		fmt.Printf("  エラー内容: %s\n", string(audioQuery))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("音声クエリの生成に失敗しました (status: %d): %s",
				resp.StatusCode, string(audioQuery)),
		})
	}

	// JSONバリデーション
	var audioQueryJSON map[string]interface{}
	if err := json.Unmarshal(audioQuery, &audioQueryJSON); err != nil {
		fmt.Printf("ERROR: 不正なJSON: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "VOICEVOXからの応答が不正です",
		})
	}

	fmt.Printf("DEBUG: audio_query successful, JSON keys: %v\n", getKeys(audioQueryJSON))

	// synthesis で音声を生成
	synthesisParams := url.Values{}
	synthesisParams.Set("speaker", fmt.Sprintf("%d", speakerID))
	synthesisURL := fmt.Sprintf("%s/synthesis?%s", voicevoxBaseURL, synthesisParams.Encode())
	fmt.Printf("リクエストURL: %s\n", synthesisURL)

	synthesisStartTime := time.Now()
	resp2, err := http.Post(synthesisURL, "application/json", bytes.NewReader(audioQuery))
	synthesisDuration := time.Since(synthesisStartTime)

	if err != nil {
		fmt.Printf("ERROR: synthesisリクエスト失敗 (%v): %v\n", synthesisDuration, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "音声合成に失敗しました: " + err.Error(),
		})
	}
	defer resp2.Body.Close()

	// ステータスコードを確認
	if resp2.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp2.Body)
		fmt.Printf("ERROR: synthesis失敗 (status: %d)\n", resp2.StatusCode)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("音声合成に失敗しました (status: %d): %s",
				resp2.StatusCode, string(bodyBytes)),
		})
	}

	// 音声データを読み込み
	audioData, err := io.ReadAll(resp2.Body)
	if err != nil {
		fmt.Printf("ERROR: 音声データ読み取り失敗: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "音声データの読み取りに失敗しました",
		})
	}

	fmt.Printf("✓ 音声データ取得完了: %d バイト (%.2f KB)\n", len(audioData), float64(len(audioData))/1024)

	totalDuration := time.Since(startTime)
	fmt.Printf("\n========================================\n")
	fmt.Printf("音声生成成功！\n")
	fmt.Printf("総処理時間: %v\n", totalDuration)
	fmt.Printf("========================================\n\n")

	// WAVファイルとして返す
	c.Response().Header().Set("Content-Type", "audio/wav")
	c.Response().Header().Set("Content-Disposition", "inline; filename=zundamon.wav")
	return c.Blob(http.StatusOK, "audio/wav", audioData)
}

func getKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
