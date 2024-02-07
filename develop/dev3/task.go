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

type ConsoleArgs struct {
	FlgColl map[string]interface{}
}

func (ca *ConsoleArgs) ParseFlags() {
	k := flag.Int("k", 0, "sort via a key; KEYDEF gives location and type")
	n := flag.Bool("n", false, "compare according to string numerical value")
	r := flag.Bool("r", false, "reverse the result of comparisons")
	u := flag.Bool("u", false, "with -c, check for strict ordering; without "+
		"-c, output only the first of an equal run")
	m := flag.Bool("M", false, "compare (unknown) < 'JAN' < ... < 'DEC'")
	b := flag.Bool("b", false, "ignore leading blanks")
	c := flag.Bool("c", false, "check for sorted input; do not sort")
	h := flag.Bool("h", false, " compare human readable numbers (e.g., 2K 1G)")

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

func (fs *FileSorter) SetFileNames(flsPth []string) {
	fs.filesPath = flsPth
}

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

func (fs *FileSorter) Sort() error {
	for _, path := range fs.filesPath {
		strVl, errRF := fs.ReadFile(path)
		if errRF != nil {
			return errRF
		}

		if fs.Flags.FlgColl["-c"].(bool) {
			item := fs.isSorted(strVl)

			if item != "" {
				fmt.Printf("disorder: %s\n", item)
			}
			return nil
		}

		fs.sorting(strVl)

		if fs.Flags.FlgColl["-r"].(bool) {
			reverse(strVl)
		}

		if fs.Flags.FlgColl["-u"].(bool) {
			strVl = unique(strVl)
		}

		strRes := fs.stringRes(strVl)

		errWrToFile := fs.writeToFile(strRes, path)

		if errWrToFile != nil {
			return errWrToFile
		}
	}
	return nil
}

func (fs *FileSorter) writeToFile(str string, filePath string) error {
	partsPath := strings.Split(filePath, "/")
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

func (fs *FileSorter) stringRes(res interface{}) string {
	resSlice, ok := res.([]string)
	if !ok {
		log.Fatal("incorrect type result")
		return ""
	}
	var strRes strings.Builder

	for i := 0; i < len(resSlice); i++ {
		strRes.WriteString(resSlice[i])
		if i < len(resSlice)-1 {
			strRes.WriteString("\n")
		}
	}

	return strRes.String()
}

func (fs *FileSorter) sorting(lines []string) {
	sort.SliceStable(lines, func(i, j int) bool {
		a := strings.Fields(lines[i])
		b := strings.Fields(lines[j])

		if fs.Flags.FlgColl["-h"].(bool) {
			return fs.compareNumericSuffix(a, b)
		}

		if fs.Flags.FlgColl["-M"].(bool) {
			return fs.compareMonths(a, b)
		}

		return fs.compareLines(a, b)
	})
}

func (fs *FileSorter) compareMonths(a, b []string) bool {
	column := fs.columnNumber()

	if column > len(a)-1 || column > len(b)-1 {
		column = 0
	}
	formatDate := "2006-January-02"

	dateA, errA := time.Parse(formatDate, fmt.Sprintf(fmt.Sprintf("2006-%s-01", a[column])))
	dateB, errB := time.Parse(formatDate, fmt.Sprintf(fmt.Sprintf("2006-%s-01", b[column])))

	if errA != nil || errB != nil {
		fmt.Println("incorrect month name")
		return false
	}

	return dateA.Before(dateB)
}

func (fs *FileSorter) compareLines(a, b []string) bool {
	column := fs.columnNumber()

	if column > len(a)-1 || column > len(b)-1 {
		column = 0
	}

	if fs.Flags.FlgColl["-b"].(bool) {
		a[column] = strings.TrimRight(a[column], " ")
		b[column] = strings.TrimRight(b[column], " ")
	}

	if fs.Flags.FlgColl["-n"].(bool) {
		numA, errA := strconv.Atoi(a[column])
		numB, errB := strconv.Atoi(b[column])

		if errA == nil && errB == nil {
			return numA < numB
		}
	}

	return a[column] < b[column]
}

func (fs *FileSorter) compareNumericSuffix(a, b []string) bool {
	column := fs.columnNumber()

	if column > len(a)-1 || column > len(b)-1 {
		column = 0
	}

	numA, sufA := fs.numericSuffix(a[column])
	numB, sufB := fs.numericSuffix(b[column])

	if numA == numB {
		return sufA < sufB
	}

	return numA < numB
}

func (fs *FileSorter) numericSuffix(s string) (int, string) {
	for i := len(s) - 1; i >= 0; i-- {
		if isNumeric(s[i]) {
			num, _ := strconv.Atoi(s[i+1:])
			return num, s[:i+1]
		}
	}

	num, _ := strconv.Atoi(s)
	return num, ""
}

func (fs *FileSorter) columnNumber() int {
	if fs.Flags.FlgColl["-k"].(int) != 0 {
		return fs.Flags.FlgColl["-k"].(int) - 1
	}

	return 0
}

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
