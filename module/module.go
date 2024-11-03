package module

import "strings"

// 受け取った文字列をぱ行に変換
func ConsonantLockLanguage(inputString string) string {

	// 1. 文字列を分割
	hiraganaCharacters := StringToSlice(inputString)

	// 2. ローマ字に変換
	romajiCharacters := HiraganaToRomaji(hiraganaCharacters)

	// 3. 母音だけ取り出し
	vowelCharacters := ExtractVowels(romajiCharacters)

	// 4. 母音をもとにぱ行に変換
	paCharacters := ConvertToHiraganaSlice(vowelCharacters)

	// 5. スライスの文字を結合
	resultString := JoinStrings(paCharacters)

	return resultString
}


// 文字列を一文字ずつ分解してスライスを返す関数
// Parameters:
//   - s: 分解対象の文字列。
// Returns:
//   - []string: 各文字が1要素として格納されたスライス。
func StringToSlice(s string) []string {
	var result []string
	for _, r := range s {
			result = append(result, string(r))
	}
	return result
}

// ひらがなをローマ字に変換
// Parameters:
//   - s: 変換対象の文字列スライス。
// Returns:
//   - []string: ローマ字に変換されたスライス
func HiraganaToRomaji(hiragana []string) []string {
romajiMap := map[string]string{
	"あ": "a", "い": "i", "う": "u", "え": "e", "お": "o",
	"か": "ka", "き": "ki", "く": "ku", "け": "ke", "こ": "ko",
	"さ": "sa", "し": "shi", "す": "su", "せ": "se", "そ": "so",
	"た": "ta", "ち": "chi", "つ": "tsu", "て": "te", "と": "to",
	"な": "na", "に": "ni", "ぬ": "nu", "ね": "ne", "の": "no",
	"は": "ha", "ひ": "hi", "ふ": "fu", "へ": "he", "ほ": "ho",
	"ま": "ma", "み": "mi", "む": "mu", "め": "me", "も": "mo",
	"や": "ya", "ゆ": "yu", "よ": "yo",
	"ら": "ra", "り": "ri", "る": "ru", "れ": "re", "ろ": "ro",
	"わ": "wa", "を": "wo", "ん": "n",
	"が": "ga", "ぎ": "gi", "ぐ": "gu", "げ": "ge", "ご": "go",
	"ざ": "za", "じ": "zi", "ず": "zu", "ぜ": "ze", "ぞ": "zo",
	"だ": "da", "ぢ": "di", "づ": "du", "で": "de", "ど": "do",
	"ば": "ba", "び": "bi", "ぶ": "bu", "べ": "be", "ぼ": "bo",
	"ぱ": "pa", "ぴ": "pi", "ぷ": "pu", "ぺ": "pe", "ぽ": "po",
}

var romaji []string
for _, char := range hiragana {
	if romajiChar, exists := romajiMap[char]; exists {
		romaji = append(romaji, romajiChar)
	} else {
		romaji = append(romaji, char)
	}
}

return romaji
}

// スライス内の文字列から母音を取り出す
// Parameters:
//   - strs: 変換対象の文字列スライス。
// Returns:
//   - []string: 各文字列から抽出された母音のスライス
func ExtractVowels(strs []string) []string {
	var result []string

	for _, s := range strs {
			var char string
			char = string(s[len(s)-1])
			result = append(result, char)
	}
	return result
}


// スライス内の文字列をひらがなに変換する関数
// Parameters:
//   - s: 変換対象の文字列スライス。
// Returns:
//   - []string: 各文字がひらがなに変換されたスライス
func ConvertToHiraganaSlice(strs []string) []string {
	var result []string

	// 母音をパ行に変換するマップ
	hiraganaMap := map[string]string{
			"a": "ぱ",
			"i": "ぴ",
			"u": "ぷ",
			"e": "ぺ",
			"o": "ぽ",
			"n": "ん",
	}

	for _, s := range strs {
			if hiraganaChar, exists := hiraganaMap[s]; exists {
					result = append(result, hiraganaChar)
			} else {
					result = append(result, "") // マッチしない場合に空文字列を追加
			}
	}
	return result
}

// 文字列のスライスを結合して1つの文字列を返す関数
// Parameters:
//   - strs: 結合する文字列のスライス
// Returns:
//   - string: 結合された文字列
func JoinStrings(strs []string) string {
	return strings.Join(strs, "")
}