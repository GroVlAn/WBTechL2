package main

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

// GroupAnagram - функция для группировки анаграмм
func GroupAnagram(strs []string) map[string][]string {
	// мапа где ключ это отсортированые символы в строке, а значение сама строка
	m := make(map[string][]string)
	for _, str := range strs {
		b := []byte(strings.ToLower(str))
		sort.Slice(b, func(i, j int) bool {
			return b[i] < b[j]
		})

		if !slices.Contains(m[string(b)], str) {
			m[string(b)] = append(m[string(b)], strings.ToLower(str))
		}
	}

	// мапа где ключ это первый элемент множества анаграмм, а значение слайс анаграм
	res := make(map[string][]string, len(m))

	for _, val := range m {
		if len(val) > 1 {
			key := val[0]
			slices.Sort(val)
			res[key] = val
		}
	}

	return res
}

func main() {
	fmt.Println(GroupAnagram([]string{
		"Тяпка",
		"пятак",
		"пятак",
		"пя1так",
		"пятак",
		"пятка",
		"тяпкА",
		"тяпкА",
		"тяпка",
		"тяпка",
		"тяПка",
		"лИСток",
		"слиток",
		"столик",
		"стол",
		"кит",
	}))
}
