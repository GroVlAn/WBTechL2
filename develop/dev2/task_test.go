package main

import "testing"

func passedString(t *testing.T, str string, correctStr string) {
	text, err := unpackString(str)
	if err != nil {
		t.Errorf("incorrect error return by string: %s", str)
	}
	if text != correctStr {
		t.Errorf("incorrect result (expected %s)", correctStr)
	}
}

func incorrectString(t *testing.T, str string) {
	text, err := unpackString(str)
	if err == nil {
		t.Errorf("incorrect error return by string: %s, expected error", str)
	}
	if text != "" {
		t.Error("incorrect result, expected empty string")
	}
}

func TestUnpacking(t *testing.T) {
	passedString(t, "a4bc2d5e", "aaaabccddddde")
	passedString(t, "abcd", "abcd")
	passedString(t, "-5", "-----")
	passedString(t, "", "")
	passedString(t, "qwe\\4\\5", "qwe45")
	passedString(t, "qwe\\45", "qwe44444")
	passedString(t, "qwe\\\\5", "qwe\\\\\\\\\\")
	passedString(t, "\\", "")
	passedString(t, "\\\\", "\\")
	passedString(t, "\\4", "4")
	incorrectString(t, "45")
	incorrectString(t, "5")
	incorrectString(t, "a45")
}
