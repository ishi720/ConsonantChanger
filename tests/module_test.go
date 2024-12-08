package tests

import (
	"myapp/module"
	"testing"
)

func TestConsonantLockLanguage(t *testing.T) {

	result1 := module.ConsonantLockLanguage("こんにちは", "pa")
	expected1 := "ぽんぴぴぱ"

	if result1 != expected1 {
		t.Errorf("文字列が一致しません。Result=%s, Expected=%s", result1, expected1)
	}

	result2 := module.ConsonantLockLanguage("ありがとう", "sa")
	expected2 := "さしさそす"
	if result2 != expected2 {
		t.Errorf("文字列が一致しません。Result=%s, Expected=%s", result2, expected2)
	}

	result3 := module.ConsonantLockLanguage("おやすみなさい", "ga")
	expected3 := "ごがぐぎががぎ"
	if result3 != expected3 {
		t.Errorf("文字列が一致しません。Result=%s, Expected=%s", result3, expected3)
	}

	// 記号のテスト
	result4 := module.ConsonantLockLanguage("わたしは、おにぎりがたべたいです。", "ma")
	expected4 := "ままみま、もみみみままめまみめむ。"
	if result4 != expected4 {
		t.Errorf("文字列が一致しません。Result=%s, Expected=%s", result4, expected4)
	}

	// 「っ」のテスト
	result5 := module.ConsonantLockLanguage("こっぺぱん", "ya")
	expected5 := "よっえやん"
	if result5 != expected5 {
		t.Errorf("文字列が一致しません。Result=%s, Expected=%s", result5, expected5)
	}

	// カタカナありのテスト
	result6 := module.ConsonantLockLanguage("ドラえもん", "sa")
	expected6 := "そさせそん"
	if result6 != expected6 {
		t.Errorf("文字列が一致しません。Result=%s, Expected=%s", result6, expected6)
	}

}
