package main

import (
	"reflect"
	"testing"
)

func TestGroupAnagram(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want map[string][]string
	}{
		{
			name: "strings has anagrams and has uppercase letters",
			args: []string{"Пятак", "пЯтка", "тяпкА", "Листок", "слИток", "стОлик", "одувАн"},
			want: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
			},
		},
		{
			name: "string without anagrams",
			args: []string{"Anagram", "Arc", "Love"},
			want: map[string][]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GroupAnagram(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupAnagram() = %v, want %v", got, tt.want)
			}
		})
	}
}
