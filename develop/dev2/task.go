package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

var ErrorIncorrectString = errors.New("incorrect string")

func addRns(strBldr *strings.Builder, run rune, count int) {
	for j := 1; j < count; j++ {
		strBldr.WriteString(string(run))
	}
}

func unpackString(str string) (string, error) {
	var res strings.Builder
	rns := []rune(str)

	for i := 0; i < len(rns); i++ {
		if i == 0 && unicode.IsDigit(rns[i]) {
			return "", ErrorIncorrectString
		}
		if i > 1 && rns[i-2] == '\\' && unicode.IsDigit(rns[i]) {
			if rns[i-1] == '\\' {
				addRns(&res, rns[i-1], int(rns[i]-'0'))
				continue
			}
			if unicode.IsDigit(rns[i-1]) {
				addRns(&res, rns[i-1], int(rns[i]-'0'))
				continue
			}
		}
		if rns[i] == '\\' {
			if i > 0 && rns[i-1] == '\\' {
				res.WriteString(string(rns[i]))
				continue
			}
			continue
		}
		if unicode.IsDigit(rns[i]) {
			if rns[i-1] == '\\' {
				res.WriteString(string(rns[i]))
				continue
			}
			if unicode.IsDigit(rns[i-1]) {
				return "", ErrorIncorrectString
			} else {
				addRns(&res, rns[i-1], int(rns[i]-'0'))
				continue
			}
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

	umpck, err := unpackString(text)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("result: %s", umpck)
}
