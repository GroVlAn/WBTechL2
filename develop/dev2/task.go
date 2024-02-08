package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

/*
ErrorIncorrectString - ошибка невено введёной строки
*/
var ErrorIncorrectString = errors.New("incorrect string")

/*
addRns - функция для добавления рун в результирующую строку
strBldr *strings.Builder - указатель на результирующую строку
rn rune - руна которую нужно записать в строку
count int - количество раз, сколько нужно добавить руну в строку
*/
func addRns(strBldr *strings.Builder, rn rune, count int) {
	for j := 1; j < count; j++ {
		strBldr.WriteRune(rn)
	}
}

/*
UnpackString - функция для распаковки строки
str string - искомая строка
*/
func UnpackString(str string) (string, error) {
	var res strings.Builder
	rns := []rune(str)

	for i := 0; i < len(rns); i++ {
		// если первая руна в строке число, возвращаем ошибку
		if i == 0 && unicode.IsDigit(rns[i]) {
			return "", ErrorIncorrectString
		}
		// проверка на \ и что текущая руна число
		if i > 1 && rns[i-2] == '\\' && unicode.IsDigit(rns[i]) {
			// если встретилось \\ записываем предыдущий \ n раз
			if rns[i-1] == '\\' {
				fmt.Println(rns[i] - '0')
				addRns(&res, rns[i-1], int(rns[i]-'0'))
				continue
			}
			// если предудущая руна число и до него есть \ и текущая руна число, выводим предыдущее число n раз
			if unicode.IsDigit(rns[i-1]) {
				addRns(&res, rns[i-1], int(rns[i]-'0'))
				continue
			}
		}

		// если текущая руна \ и предыдущие условия не выполнились
		if rns[i] == '\\' {
			// если это не первый элемент и предыдущий равен \ то просто выводим \
			if i > 0 && rns[i-1] == '\\' {
				res.WriteString(string(rns[i]))
				continue
			}
			continue
		}
		// если текущая руна число и предудущие условия не выполнились
		if unicode.IsDigit(rns[i]) {
			// если предвдущая руна \ то просто выводим число
			if rns[i-1] == '\\' {
				res.WriteString(string(rns[i]))
				continue
			}
			// если предыдущая руна тоже число, значит строка некорректна иначе выводим предыдущую руну n раз
			if unicode.IsDigit(rns[i-1]) {
				return "", ErrorIncorrectString
			} else {
				addRns(&res, rns[i-1], int(rns[i]-'0'))
				continue
			}
			// иначе просто записываем руну в результат
		} else {
			res.WriteString(string(rns[i]))
			continue
		}
	}

	return res.String(), nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter string ->")
	text, errRead := reader.ReadString('\n')

	if errRead != nil {
		fmt.Println(errRead)
		return
	}

	umpck, err := UnpackString(text)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("result: %s", umpck)
}
