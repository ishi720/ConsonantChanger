async function convertButton() {
    const inputString = document.getElementById("input_string").value;
    const lineType = document.getElementById("line_type").value;

    try {
        // APIリクエストのURLを作成
        const response = await fetch(`/api/getColockLanguage?input_string=${encodeURIComponent(inputString)}&line_type=${encodeURIComponent(lineType)}`);
        
        // エラーチェック
        if (!response.ok) {
          throw new Error('APIリクエストに失敗しました');
        }

        // レスポンスをJSONとして解析
        const data = await response.json();

        // 結果を表示
        document.getElementById("result").textContent = data.result;
    } catch (error) {
        console.error("エラー:", error);
        document.getElementById("result").textContent = "エラーが発生しました";
    }
}