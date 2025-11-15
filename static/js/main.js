document.addEventListener('DOMContentLoaded', function() {
    var elems = document.querySelectorAll('select');
    var instances = M.FormSelect.init(elems);

    // 入力時に変換ボタンの有効化
    var inputString = document.getElementById("input_string");
    inputString.addEventListener("input", function() {
        document.getElementById("convert_btn").classList.toggle("disabled", inputString.value.length === 0);
    });

    // コピーボタンのイベント
    document.getElementById('copyButton').addEventListener('click', () => {
        let textContent = document.getElementById("result").value;
        navigator.clipboard.writeText(textContent)
            .then(() => {
                console.log('クリップボードにコピーしました:', textContent);
                M.toast({html: 'コピーしました！'});
            })
            .catch(err => {
                console.error('コピーに失敗しました:', err);
                M.toast({html: 'コピーに失敗しました'});
            });
    });

    // 音声再生ボタンのイベント
    document.getElementById('playButton').addEventListener('click', playVoice);
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

        // コピーボタンと音声再生ボタンを有効化
        document.getElementById('copyButton').removeAttribute('hidden');
        document.getElementById('playButton').removeAttribute('hidden');
    } catch (error) {
        console.error("エラー:", error);
        document.getElementById("result").textContent = "エラーが発生しました";
        M.toast({html: 'エラーが発生しました'});
    }
}

async function playVoice() {
    const resultText = document.getElementById("result").value;
    
    if (!resultText) {
        M.toast({html: '再生するテキストがありません'});
        return;
    }

    const playButton = document.getElementById('playButton');
    const originalHTML = playButton.innerHTML;
    
    try {
        // ボタンを無効化
        playButton.classList.add('disabled');
        playButton.innerHTML = '<i class="material-icons left">hourglass_empty</i>生成中...';

        console.log('音声生成リクエスト:', resultText);

        // POSTリクエストで音声生成APIを呼び出し
        const response = await fetch('/api/generateVoice', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                text: resultText
            })
        });
        
        console.log('レスポンスステータス:', response.status);
        console.log('Content-Type:', response.headers.get('Content-Type'));

        if (!response.ok) {
            // エラーレスポンスの詳細を取得
            const contentType = response.headers.get('Content-Type');
            let errorMessage = '音声生成に失敗しました';
            
            if (contentType && contentType.includes('application/json')) {
                const errorData = await response.json();
                errorMessage = errorData.error || errorMessage;
                console.error('エラー詳細:', errorData);
            } else {
                const errorText = await response.text();
                console.error('エラーレスポンス:', errorText);
                errorMessage = errorText || errorMessage;
            }
            
            throw new Error(errorMessage);
        }

        // 音声データを取得
        const audioBlob = await response.blob();
        console.log('音声データサイズ:', audioBlob.size, 'bytes');
        
        if (audioBlob.size < 100) {
            throw new Error('音声データが小さすぎます（生成に失敗した可能性があります）');
        }
        
        const audioUrl = URL.createObjectURL(audioBlob);

        // 音声を再生
        const audio = new Audio(audioUrl);
        
        audio.onended = () => {
            URL.revokeObjectURL(audioUrl);
            playButton.classList.remove('disabled');
            playButton.innerHTML = originalHTML;
        };

        audio.onerror = (e) => {
            console.error('音声再生エラー:', e);
            URL.revokeObjectURL(audioUrl);
            playButton.classList.remove('disabled');
            playButton.innerHTML = originalHTML;
            M.toast({html: '音声の再生に失敗しました'});
        };

        console.log('音声再生開始');
        await audio.play();
        M.toast({html: 'ずんだもんが読み上げます！'});

    } catch (error) {
        console.error("音声再生エラー:", error);
        M.toast({html: error.message || '音声生成に失敗しました。VOICEVOXが起動しているか確認してください。'});
        playButton.classList.remove('disabled');
        playButton.innerHTML = originalHTML;
    }
}