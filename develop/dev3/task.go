package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var suffixes = map[string]float64{
	"K": 1e3,
	"M": 1e6,
	"G": 1e9,
	"T": 1e12,
	"P": 1e15,
	"E": 1e18,
	"Z": 1e21,
	"Y": 1e24,
}

/*
ConsoleArgs - структура хранящая доступные флаги
FlgColl map[string]interface{} - мапа хранаящая имя флага: значение
*/
type ConsoleArgs struct {
	FlgColl map[string]interface{}
}

// ParseFlags - парсим флаги из консоли
func (ca *ConsoleArgs) ParseFlags() {
	k := flag.Int("k", 0, "sort via a key; KEYDEF gives location and type")
	n := flag.Bool("n", false, "compare according to string numerical value")
	r := flag.Bool("r", false, "reverse the result of comparisons")
	u := flag.Bool("u", false, "with -c, check for strict ordering; without "+
		"-c, output only the first of an equal run")
	m := flag.Bool("M", false, "compare (unknown) < 'JAN' < ... < 'DEC'")
	b := flag.Bool("b", false, "ignore leading blanks")
	c := flag.Bool("c", false, "check for sorted input; do not sort")
	h := flag.Bool("h", false, "compare human readable numbers (e.g., 2K 1G)")

	flag.Parse()

	ca.FlgColl = make(map[string]interface{}, 8)
	ca.FlgColl["-k"] = *k
	ca.FlgColl["-n"] = *n
	ca.FlgColl["-r"] = *r
	ca.FlgColl["-u"] = *u
	ca.FlgColl["-M"] = *m
	ca.FlgColl["-b"] = *b
	ca.FlgColl["-c"] = *c
	ca.FlgColl["-h"] = *h
}

// ReadFilePaths - последними аргументами ожидаем передачи пути до файлов, в которых нужно отсортировать данные
func (ca *ConsoleArgs) ReadFilePaths() []string {
	args := os.Args
	filesPath := make([]string, 0)

	for i := len(args) - 1; i >= 1; i-- {
		_, ok := ca.FlgColl[args[i]]

		if ok || (!strings.Contains(args[i], "\"") && !strings.Contains(args[i], "/")) {
			break
		}
		filesPath = append(filesPath, args[i])
	}
	return filesPath
}

/*
FileSorter - структура занимающаяся сортировкой файла
Flags *ConsoleArgs - указатель на структуру флагов
filesPath []string - пути до файлов
*/
type FileSorter struct {
	Flags     *ConsoleArgs
	filesPath []string
}

func NewFileSorter(flgs *ConsoleArgs) *FileSorter {
	return &FileSorter{
		Flags:     flgs,
		filesPath: make([]string, 0),
	}
}

// SetFileNames - метод для инициализации путей к файлам
func (fs *FileSorter) SetFileNames(flsPth []string) {
	fs.filesPath = flsPth
}

/*
ReadFile - метод читающий строки из файла
filePath string - путь до файла
*/
func (fs *FileSorter) ReadFile(filePath string) ([]string, error) {
	file, errOpen := os.Open(filePath)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Can not close file: %s", err.Error())
		}
	}(file)
	strValues := make([]string, 0)

	if errOpen != nil {
		return nil, fmt.Errorf("file not exist: %s", errOpen.Error())
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	success := scanner.Scan()
	for success {
		if !success {
			err := scanner.Err()

			if err != nil {
				return nil, fmt.Errorf("can not read file: %s", err.Error())
			}
			break
		}

		strValues = append(strValues, scanner.Text())
		success = scanner.Scan()
	}

	return strValues, nil
}

// Sort - метод для сортировки файла
func (fs *FileSorter) Sort() error {
	for _, path := range fs.filesPath {
		strVl, errRF := fs.ReadFile(path)
		if errRF != nil {
			return errRF
		}

		// проверка на отсортированость файла, если установлен соответствующий флаг
		if fs.Flags.FlgColl["-c"].(bool) {
			item := fs.isSorted(strVl)

			if item != "" {
				fmt.Printf("disorder: %s\n", item)
			}
			return nil
		}

		// сортируем файл
		fs.sorting(strVl)

		// если установлен флаг разворачиваем результат
		if fs.Flags.FlgColl["-r"].(bool) {
			reverse(strVl)
		}

		// если установлен флаг, оставляем только уникальные значения
		if fs.Flags.FlgColl["-u"].(bool) {
			strVl = unique(strVl)
		}

		// превращаем массив строк в единую строку
		strRes := fs.stringRes(strVl)

		// записываем результат в файл
		errWrToFile := fs.writeToFile(strRes, path)

		if errWrToFile != nil {
			return errWrToFile
		}
	}
	return nil
}

/*
writeToFile - метод для записи результата в файл
str string - результат, который нужно записать в файл
filePath string - путь до текущего файла
*/
func (fs *FileSorter) writeToFile(str string, filePath string) error {
	partsPath := strings.Split(filePath, "/")
	//создаём новый файл с приставкой _res
	nameFile := partsPath[len(partsPath)-1] + "_res"
	file, errOpen := os.OpenFile(
		nameFile,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)

	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Fatalf("cat not close file: %s", err.Error())
		}
	}(file)

	if errOpen != nil {
		return fmt.Errorf("can not open file: %s", errOpen.Error())
	}

	_, errWrite := file.Write([]byte(str))

	if errWrite != nil {
		return fmt.Errorf("can not write to file: %s", errWrite.Error())
	}

	return nil
}

/*
stringRes - метод для превразения массива строк в строку
res []string - массив строк
*/
func (fs *FileSorter) stringRes(res []string) string {
	var strRes strings.Builder
	resSlice := res

	for i := 0; i < len(resSlice); i++ {
		strRes.WriteString(resSlice[i])
		if i < len(resSlice)-1 {
			strRes.WriteString("\n")
		}
	}

	return strRes.String()
}

/*
sorting - метод для сортировки строк
*/
func (fs *FileSorter) sorting(lines []string) {
	sort.SliceStable(lines, func(i, j int) bool {
		a := strings.Fields(lines[i])
		b := strings.Fields(lines[j])

		// если указан флаг в сравниваем с учётом суффикса
		if fs.Flags.FlgColl["-h"].(bool) {
			return fs.compareNumericSuffix(a, b)
		}

		// если установлен флаг сравниваем месяцы
		if fs.Flags.FlgColl["-M"].(bool) {
			return fs.compareMonths(a, b)
		}

		// иначе сравниваем строки как есть
		return fs.compareLines(a, b)
	})
}

// compareMonths - метод для сравнения месяцев
func (fs *FileSorter) compareMonths(a, b []string) bool {
	column := fs.columnNumber()

	if column > len(a)-1 || column > len(b)-1 {
		column = 0
	}
	formatDate := "2006-January-02"

	// приводим значения к дате
	dateA, errA := time.Parse(formatDate, fmt.Sprintf("2006-%s-01", a[column]))
	dateB, errB := time.Parse(formatDate, fmt.Sprintf("2006-%s-01", b[column]))

	if errA != nil || errB != nil {
		fmt.Println("incorrect month name")
		return false
	}

	return dateA.Before(dateB)
}

// compareLines - метод для сравнения строк
func (fs *FileSorter) compareLines(a, b []string) bool {
	column := fs.columnNumber()

	if column > len(a)-1 || column > len(b)-1 {
		column = 0
	}

	// если установлен флаг, игнорируем хвостовые пробелы
	if fs.Flags.FlgColl["-b"].(bool) {
		a[column] = strings.TrimRight(a[column], " ")
		b[column] = strings.TrimRight(b[column], " ")
	}

	// если установлен флаг, сравниваем строку как число
	if fs.Flags.FlgColl["-n"].(bool) {
		numA, errA := strconv.Atoi(a[column])
		numB, errB := strconv.Atoi(b[column])

		if errA == nil && errB == nil {
			return numA < numB
		}
	}

	// иначе сравниваем символы определёной колонки как строки
	return a[column] < b[column]
}

// compareNumericSuffix - метод для сравнения символов с суффиксом
func (fs *FileSorter) compareNumericSuffix(a, b []string) bool {
	column := fs.columnNumber()

	if column > len(a)-1 || column > len(b)-1 {
		column = 0
	}

	// приводим числа с суфиксами к числу
	numA, errA := fs.numericSuffix(a[column])
	numB, errB := fs.numericSuffix(b[column])

	// ошибки при конвертации нет, сравниваем полученные числа
	if errA == nil && errB == nil {
		return numA < numB
	}

	// иначе просто сравниваем как строки
	return a[column] < b[column]
}

// numericSuffix - метод приводящий число с суффиксом к int64
func (fs *FileSorter) numericSuffix(s string) (int64, error) {
	for suf, mul := range suffixes {
		if strings.HasSuffix(s, suf) {
			strSuf := strings.TrimSuffix(s, suf)
			val, err := strconv.ParseInt(strSuf, 10, 64)

			if err != nil {
				return 0, err
			}

			return int64(mul) * val, nil
		}
	}

	return strconv.ParseInt(s, 10, 64)
}

// columnNumber - метод для получения номера колонки
func (fs *FileSorter) columnNumber() int {
	if fs.Flags.FlgColl["-k"].(int) != 0 {
		return fs.Flags.FlgColl["-k"].(int) - 1
	}

	return 0
}

// isSorted - метод для поверки сотсортированности файла
func (fs *FileSorter) isSorted(lines []string) string {
	var cur string

	for i := 0; i < len(lines); i++ {
		items := strings.Fields(lines[i])
		for j := 0; j < len(items); j++ {
			if i == 0 && j == 0 {
				cur = items[j]
				continue
			}
			if items[j] < cur && !fs.Flags.FlgColl["-r"].(bool) {
				return items[j]
			} else if items[j] > cur && fs.Flags.FlgColl["-r"].(bool) {
				return items[j]
			}
		}
	}

	return ""
}

func main() {
	clsArgs := &ConsoleArgs{}
	clsArgs.ParseFlags()
	fileSorter := NewFileSorter(clsArgs)
	fileSorter.SetFileNames(clsArgs.ReadFilePaths())
	err := fileSorter.Sort()

	if err != nil {
		fmt.Println(err)
	}
}

func reverse(strs []string) {
	for i, j := 0, len(strs)-1; i < j; i, j = i+1, j-1 {
		strs[i], strs[j] = strs[j], strs[i]
	}
}

func unique(strs []string) []string {
	keys := make(map[string]bool, len(strs))
	res := make([]string, 0)

	for _, entry := range strs {
		if _, ok := keys[entry]; !ok {
			keys[entry] = true
			res = append(res, entry)
		}
	}

	return res
}

func isNumeric(c byte) bool {
	return c >= '0' && c <= '9'
}
