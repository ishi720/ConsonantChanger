package module

import "strings"

// 受け取った文字列をぱ行に変換
func ConsonantLockLanguage(inputString string, lineType string) string {

	// 1. 文字列を分割
	hiraganaCharacters := StringToSlice(inputString)

	// 2. ローマ字に変換
	romajiCharacters := HiraganaToRomaji(hiraganaCharacters)

	// 3. 母音だけ取り出し
	vowelCharacters := ExtractVowels(romajiCharacters)

	// 4. 母音をもとにぱ行に変換
	paCharacters := ConvertToHiraganaSlice(vowelCharacters, lineType)

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

// スライス内の文字列をひらがなに変換する
// Parameters:
//   - strs: 変換対象の文字列スライス。
//   - lineType: ひらがなの行。
// Returns:
//   - []string: 各文字がひらがなに変換されたスライス。一致する行がなければ空で返す。
func ConvertToHiraganaSlice(strs []string, lineType string) []string {
	var result []string

	// 母音をパ行に変換するマップ
	lineMaps := map[string]map[string]string{
		"a":  {"a": "あ", "i": "い", "u": "う", "e": "え", "o": "お", "n": "ん"},
		"ka": {"a": "か", "i": "き", "u": "く", "e": "け", "o": "こ", "n": "ん"},
		"sa": {"a": "さ", "i": "し", "u": "す", "e": "せ", "o": "そ", "n": "ん"},
		"ta": {"a": "た", "i": "ち", "u": "つ", "e": "て", "o": "と", "n": "ん"},
		"na": {"a": "な", "i": "に", "u": "ぬ", "e": "ね", "o": "の", "n": "ん"},
		"ha": {"a": "は", "i": "ひ", "u": "ふ", "e": "へ", "o": "ほ", "n": "ん"},
		"ma": {"a": "ま", "i": "み", "u": "む", "e": "め", "o": "も", "n": "ん"},
		"ya": {"a": "や", "i": "い", "u": "ゆ", "e": "え", "o": "よ", "n": "ん"},
		"ra": {"a": "ら", "i": "り", "u": "る", "e": "れ", "o": "ろ", "n": "ん"},
		"ga": {"a": "が", "i": "ぎ", "u": "ぐ", "e": "げ", "o": "ご", "n": "ん"},
		"za": {"a": "ざ", "i": "じ", "u": "ず", "e": "ぜ", "o": "ぞ", "n": "ん"},
		"da": {"a": "だ", "i": "ぢ", "u": "づ", "e": "で", "o": "ど", "n": "ん"},
		"ba": {"a": "ば", "i": "び", "u": "ぶ", "e": "べ", "o": "ぼ", "n": "ん"},
		"pa": {"a": "ぱ", "i": "ぴ", "u": "ぷ", "e": "ぺ", "o": "ぽ", "n": "ん"},
	}

	// 指定された行に基づいて適切なマップを選択
	hiraganaMap, exists := lineMaps[lineType]
	// 行が存在しない場合、空の結果を返す
	if !exists {
		return result
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
