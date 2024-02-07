package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type ConsoleArgs struct {
	FlgColl map[string]interface{}
}

func (ca *ConsoleArgs) ParseFlags() {
	f := flag.Int("f", 0, "chose column")
	d := flag.String("d", "\t", "change delimiter")
	s := flag.Bool("s", false, "print only line with delimiter")

	flag.Parse()

	ca.FlgColl = make(map[string]interface{}, 8)
	ca.FlgColl["-f"] = *f
	ca.FlgColl["-d"] = *d
	ca.FlgColl["-s"] = *s
}

type Cut struct {
	ca       *ConsoleArgs
	filePath string
}

func NewCut(ca *ConsoleArgs) *Cut {
	return &Cut{
		ca: ca,
	}
}

func (c *Cut) Cut() ([][]string, error) {
	lines, err := c.readFile()

	if err != nil {
		return nil, err
	}

	if c.ca.FlgColl["-f"].(int) != 0 {
		numCol := c.ca.FlgColl["-f"].(int) - 1
		curRes := make([][]string, 0)

		for _, line := range lines {
			ln := make([]string, 0)
			for key, item := range line {
				if key == numCol {
					ln = append(ln, item)
				}
			}
			if len(ln) > 0 {
				curRes = append(curRes, ln)
			}
		}

		if len(curRes) > 0 {
			return curRes, nil
		}
	}

	return lines, nil
}

func (c *Cut) SetFilePath() error {
	args := os.Args

	if len(args) < 2 {
		return fmt.Errorf("file path is empty")
	}

	c.filePath = args[len(args)-1]

	return nil
}

func (c *Cut) readFile() ([][]string, error) {
	file, errFile := os.Open(c.filePath)
	defer func(file *os.File) {
		if errClose := file.Close(); errClose != nil {
			panic(errClose)
			log.Fatal("can not close file")
		}
	}(file)

	if errFile != nil {
		return nil, fmt.Errorf("file %s does not exist", c.filePath)
	}

	lines := make([][]string, 0)
	scanner := bufio.NewScanner(file)
	success := scanner.Scan()
	delimiter := c.ca.FlgColl["-d"].(string)

	for success {
		text := scanner.Text()
		if c.ca.FlgColl["-s"].(bool) && !strings.Contains(text, delimiter) {
			success = scanner.Scan()
			continue
		}
		col := strings.Split(scanner.Text(), delimiter)
		lines = append(lines, col)
		success = scanner.Scan()
	}

	return lines, nil
}

func main() {
	ca := &ConsoleArgs{}
	ca.ParseFlags()

	cut := NewCut(ca)
	errFile := cut.SetFilePath()

	if errFile != nil {
		fmt.Println(errFile.Error())
		return
	}

	res, err := cut.Cut()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, line := range res {
		for _, item := range line {
			fmt.Printf("%s ", item)
		}
		fmt.Print("\n")
	}
}
