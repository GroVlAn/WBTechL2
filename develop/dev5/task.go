package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

/*
ConsoleArgs - структура хранящая доступные флаги
FlgColl map[string]interface{} - мапа хранаящая имя флага: значение
*/
type ConsoleArgs struct {
	FlgColl map[string]interface{}
}

// CollectAfter - структура отвечающая за получение элементов после найденого
type CollectAfter struct {
	ca *ConsoleArgs
}

// CollectBefore - структура отвечающая за получение элементов до найденого
type CollectBefore struct {
	ca *ConsoleArgs
}

// CollectAB - структура отвечающая за получение элементов до и после найденого
type CollectAB struct {
	ca *ConsoleArgs
}

func NewGrep(ca *ConsoleArgs) *Grep {
	return &Grep{
		ca: ca,
		colA: &CollectAfter{
			ca: ca,
		},
		colB: &CollectBefore{
			ca: ca,
		},
		colC: &CollectAB{
			ca: ca,
		},
	}
}

/*
Grep - структура фильтрующая данные из файла
ca       *ConsoleArgs - аргументы с которыми запущена программа
colA     Collector - структура для получения элементов после найденого, реализует интерфейс Collector
colB     Collector - структура для получения элементов до найденого, реализует интерфейс Collector
colC     Collector - структура для получения элементов до и после найденого, реализует интерфейс Collector
target   string - искомое значение
filePath string - путь до файла
*/
type Grep struct {
	ca       *ConsoleArgs
	colA     Collector
	colB     Collector
	colC     Collector
	target   string
	filePath string
}

// ParseFlags - метод парсит флаги из консоли
func (ca *ConsoleArgs) ParseFlags() {
	A := flag.Int("A", 0, "print n lines after match")
	B := flag.Int("B", 0, "print n lines before match")
	C := flag.Int("C", 0, "print n lines after and before match")
	c := flag.Bool("c", false, "print count lines")
	i := flag.Bool("i", false, "ignoring register")
	v := flag.Bool("v", false, "invert filter")
	F := flag.Bool("F", false, "full line match")
	n := flag.Bool("n", false, "print number of lines match")

	flag.Parse()

	ca.FlgColl = make(map[string]interface{}, 8)
	ca.FlgColl["-A"] = *A
	ca.FlgColl["-B"] = *B
	ca.FlgColl["-C"] = *C
	ca.FlgColl["-c"] = *c
	ca.FlgColl["-i"] = *i
	ca.FlgColl["-v"] = *v
	ca.FlgColl["-F"] = *F
	ca.FlgColl["-n"] = *n
}

/*
Collect -  метод для получения значения после найденого
col []string - все строки из файла
matched []string - строки равные искомому
keys []int - номера строк в файле
*/
func (colA *CollectAfter) Collect(col []string, matched []string, keys []int) ([]string, []int) {
	after, okA := colA.ca.FlgColl["-A"].(int)

	if !okA || after <= 0 || len(matched) == 0 {
		return matched, keys
	}

	return addAfter(col, matched, keys, after)
}

/*
Collect - метод для получения значения до найденого
col []string - все строки из файла
matched []string - строки равные искомому
keys []int - номера строк в файле
*/
func (colB *CollectBefore) Collect(col []string, matched []string, keys []int) ([]string, []int) {
	before, okB := colB.ca.FlgColl["-B"].(int)

	if !okB || before <= 0 || len(matched) == 0 {
		return matched, keys
	}

	return addBefore(col, matched, keys, before)
}

/*
Collect - реализация метода получения значений до и после найденого
col []string - все строки из файла
matched []string - строки равные искомому
keys []int - номера строк в файле
*/
func (colAB *CollectAB) Collect(col []string, matched []string, keys []int) ([]string, []int) {
	ab, okB := colAB.ca.FlgColl["-C"].(int)

	if !okB || ab <= 0 || len(matched) == 0 {
		return matched, keys
	}
	res, resKeys := addBefore(col, matched, keys, ab)
	res, resKeys = addAfter(col, res, resKeys, ab)

	return res, resKeys
}

// SetArgs - читаем из консоли путь до файла и искомый элемент
func (g *Grep) SetArgs() error {
	args := os.Args

	if len(args) < 3 {
		return fmt.Errorf("file path is empty")
	}

	if len(args) < 2 {
		return fmt.Errorf("target and file path is empty")
	}

	errRFP := g.readFilePath(args)

	if errRFP != nil {
		return errRFP
	}

	g.readTarget(args)

	return nil
}

// Grep - фильтруем строки из файла
func (g *Grep) Grep() ([]string, []int, error) {
	lines, errRF := g.readFile()

	if errRF != nil {
		return nil, nil, errRF
	}

	// слайсы: результирующий и номеров строк в файле
	res := make([]string, 0, cap(lines))
	keysRes := make([]int, 0, cap(lines))

	var ignCase string

	// если установлен флаг игнорируем регистр строки
	if g.ca.FlgColl["-i"].(bool) {
		ignCase = "(?i)"
	}

	for key, line := range lines {
		matched, err := regexp.Match(ignCase+g.target, []byte(line))

		// если установлен флаг проверяем, что вся строка равна искомой
		if g.ca.FlgColl["-F"].(bool) && g.target == line {
			res = append(res, line)
			keysRes = append(keysRes, key)
		}

		// если установлен флаг сохраняем строки не удовлетворяющие искомой
		if !matched && g.ca.FlgColl["-v"].(bool) {
			res = append(res, line)
			keysRes = append(keysRes, key)
		}

		// если же флагов для условия сравнения нет, проверяем есть ли искомая строка в текущей из файла
		if err == nil && matched && !g.ca.FlgColl["-v"].(bool) && !g.ca.FlgColl["-F"].(bool) {
			res = append(res, line)
			keysRes = append(keysRes, key)
		}
	}

	// вызываем методы для получения строк до и после, сами методы решат, нужно ли добавлять значения к текущему слайсу
	res, keysRes = g.colA.Collect(lines, res, keysRes)
	res, keysRes = g.colB.Collect(lines, res, keysRes)
	res, keysRes = g.colC.Collect(lines, res, keysRes)

	return res, keysRes, nil
}

// readTarget - метод для получения искомой строки
func (g *Grep) readTarget(args []string) {
	g.target = args[len(args)-2]
}

// readFilePath - метод для получения пути до файла
func (g *Grep) readFilePath(args []string) error {
	path := args[len(args)-1]

	if !g.fileExist(path) {
		return fmt.Errorf("file %s does not exist", path)
	}

	g.filePath = args[len(args)-1]

	return nil
}

// readFile - метод для получения строк из файла
func (g *Grep) readFile() ([]string, error) {
	file, errFile := os.Open(g.filePath)

	defer func(file *os.File) {
		if errClose := file.Close(); errClose != nil {
			log.Fatalf("Can not close file")
		}
	}(file)

	if errFile != nil {
		return nil, fmt.Errorf("file %s does not exist", file.Name())
	}

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)

	success := scanner.Scan()

	for success {
		lines = append(lines, scanner.Text())
		success = scanner.Scan()
	}

	return lines, nil
}

// fileExist - метод для проверки существования файла
func (g *Grep) fileExist(path string) bool {
	_, err := os.Stat(path)

	if errors.Is(err, os.ErrNotExist) || err != nil {
		return false
	}

	return true
}

func main() {
	ca := &ConsoleArgs{}
	ca.ParseFlags()

	grep := NewGrep(ca)
	errArgs := grep.SetArgs()

	if errArgs != nil {
		fmt.Println(errArgs.Error())
		return
	}

	res, keys, errGrep := grep.Grep()

	if errGrep != nil {
		fmt.Println(errGrep.Error())
		return
	}

	if len(res) > 0 && len(keys) > 0 {
		if ca.FlgColl["-n"].(bool) {
			for _, key := range keys {
				fmt.Println(key)
			}
			return
		}
		if ca.FlgColl["-c"].(bool) {
			fmt.Println(len(res))
			return
		}
		for _, line := range res {
			fmt.Println(line)
		}
		return
	}

	fmt.Println("empty")
}

// Collector - интерфейс для структур которые могут отдавать найденые значения + новые
type Collector interface {
	Collect(col []string, matched []string, keys []int) ([]string, []int)
}

/*
addAfter - функция для получения элементов из слайса до найденого
col []string - все строки из файла
matched []string - строки равные искомому
keys []int - номера строк в файле
*/
func addAfter(col []string, matched []string, keys []int, after int) ([]string, []int) {
	lastKey := keys[len(keys)-1]

	for i := lastKey + 1; i <= lastKey+after; i++ {
		if i >= len(col) {
			break
		}
		matched = append(matched, col[i])
		keys = append(keys, i)
	}

	return matched, keys
}

/*
addBefore - функция для получения элементов из слайса до найденого
col []string - все строки из файла
matched []string - строки равные искомому
keys []int - номера строк в файле
*/
func addBefore(col []string, matched []string, keys []int, before int) ([]string, []int) {
	firstKey := keys[0]
	newRes := make([]string, 0)
	newKeys := make([]int, 0)

	for i := firstKey - before; i < firstKey; i++ {
		if i < 0 {
			break
		}
		newRes = append(newRes, col[i])
		newKeys = append(newKeys, i)
	}

	if len(newRes) > 0 {
		newRes = append(newRes, matched...)
		newKeys = append(newKeys, keys...)
		return newRes, newKeys
	}

	return matched, keys
}
