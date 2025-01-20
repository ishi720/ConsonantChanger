document.addEventListener('DOMContentLoaded', function() {
    var elems = document.querySelectorAll('select');
    var instances = M.FormSelect.init(elems);

    // 入力時に変換ボタンの有効化
    var inputString = document.getElementById("input_string");
    inputString.addEventListener("input", function() {
        document.getElementById("convert_btn").classList.toggle("disabled", inputString.value.length === 0);
    });

    // 変換ボタンのイベント
    var convert_btn = document.getElementById("convert_btn");
    var outputString = document.getElementById("result");
    convert_btn.addEventListener("click", function() {
        setTimeout(function() {
            document.getElementById('copyButton').removeAttribute('hidden');
        }, 500);
    });

    // コピーボタンのイベント
    document.getElementById('copyButton').addEventListener('click', () => {
        let textContent = document.getElementById("result").value;
        navigator.clipboard.writeText(textContent)
            .then(() => {
                console.log('クリップボードにコピーしました:', textContent);
            })
            .catch(err => {
                console.error('コピーに失敗しました:', err);
            });
    });
});

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