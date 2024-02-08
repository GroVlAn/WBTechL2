package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"
)

/*
Directories - структура отвечающая за работу с директориями ОС
curDir   strings.Builder - текущая директория
backup   string - предыдущая директория - нужна для возврата к ней, если новая не существует
homePath string - домашняя директория, к воторой вы вошли в программу
*/
type Directories struct {
	curDir   strings.Builder
	backup   string
	homePath string
}

// Stdout - структура отвечающая за вывод в STDOUT
type Stdout struct{}

// Processes - структура отвечающая за работу с процессами OC
type Processes struct{}

/*
Shell - структура - консоль
dir DirWorker - стуктура, которая может работать с директориями
std StdWorker - структура, которая может работать с выводом в консоль
prc PrcWorker - структура, которая может работать с процессами ОС
*/
type Shell struct {
	dir DirWorker
	std StdWorker
	prc PrcWorker
}

func NewDirectories(dir string) *Directories {
	directories := &Directories{}
	directories.curDir.WriteString(dir)
	directories.homePath = dir
	directories.backup = dir
	return directories
}

func NewStdout() *Stdout {
	return &Stdout{}
}

func NewProcesses() *Processes {
	return &Processes{}
}

func NewShell(dir DirWorker, std StdWorker, prc PrcWorker) *Shell {
	return &Shell{
		dir: dir,
		std: std,
		prc: prc,
	}
}

/*
Execute - метод выполняющий текущие команды
cmdArgs string - текущая строка содержащая набор команд
*/
func (sh *Shell) Execute(cmdArgs string) {

	// разделяем строку на отдельные команды
	cmds := strings.Split(cmdArgs, "|")

	for _, cmd := range cmds {
		cmd = strings.Trim(cmd, " ")
		var args []string
		lineParts := strings.Split(cmd, " ")
		cmd := lineParts[0]
		if len(lineParts) > 1 {
			args = lineParts[1:]
		}
		// выполняем команду, передав саму команду и её аргументы
		sh.executeCommand(cmd, args)
	}

}

/*
executeCommand - метод выполняющий конкретную команду
cmd string - название команды
args []string - аргументы команды
*/
func (sh *Shell) executeCommand(cmd string, args []string) {
	// определяем по названию команды, какой метод нужно вызвать
	switch cmd {
	case "cd":
		if len(args) == 0 {
			sh.dir.changeDir("")
			return
		}
		sh.dir.changeDir(args[len(args)-1])
	case "pwd":
		sh.dir.pwd()
	case "echo":
		sh.std.echo(args...)
	case "kill":
		for _, arg := range args {
			err := sh.prc.kill(arg)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	case "ps":
		err := sh.prc.ps(args)

		if err != nil {
			fmt.Println(err.Error())
		}
	default:
		err := sh.exec(cmd, args)

		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

/*
exec - метод запрашивающий у ОС команду, которую нужно выполнить
cmd string - название команды
args []string - аргументы команды
*/
func (sh *Shell) exec(cmd string, args []string) error {
	command, err := exec.Command(cmd, args...).Output()

	if err != nil {
		return fmt.Errorf("exec: %s\n", err.Error())
	}

	fmt.Print(string(command))
	return nil
}

/*
changeDir - метод меняющий текущую дерикторию
dir string - дериктория, куда нужно перейти
*/
func (d *Directories) changeDir(dir string) {
	// если dir пустой, переходим к дериктории с которой начали
	if len(dir) == 0 {
		d.curDir.Reset()
		d.curDir.WriteString(d.homePath)
		return
	}
	// разбиваем новый путь на составляющие
	parts := strings.Split(dir, "/")
	isNext := false
	isResetNotDots := false

	d.backup = d.curDir.String()

	for _, part := range parts {
		//если текущий кусок пуст, игнорируем его
		if part == "" {
			continue
		}
		curDir := d.curDir.String()
		// если кусок равен .. удаляем послуднюю часть текущей директории
		if part == ".." {
			lstSlh := strings.LastIndex(curDir, "/")

			if lstSlh == -1 {
				continue
			}
			d.curDir.Reset()
			d.curDir.WriteString(curDir[:lstSlh])
			continue
		}
		// если текущий кусок равен . значит ставим флаг что нужно добавлять к уже текущей директории
		// иначе сбрасываем текущую директорию
		if part == "." {
			isNext = true
			continue
		} else if !isResetNotDots && !isNext {
			d.curDir.Reset()
			isResetNotDots = true
		}
		// иначе прото добавляем к текущей директории новый кусок
		d.curDir.WriteString("/")
		d.curDir.WriteString(part)
	}

	// если текущего пути не существует в ОС, сбрасываем текущую директорию до предыдущей
	if !d.exist() {
		fmt.Printf("cd: path %s does not exist\n", d.curDir.String())
		d.curDir.Reset()
		d.curDir.WriteString(d.backup)
	}
}

// exist - метод проверяющий директорию на существование
func (d *Directories) exist() bool {
	_, err := os.Stat(d.curDir.String())

	if err != nil {
		return false
	}

	return true
}

// pwd - метод выводящий в консоль текущую директорию
func (d *Directories) pwd() {
	fmt.Println(d.curDir.String())
}

// echo - метод выводящий в консоль
func (std *Stdout) echo(strs ...string) {
	if len(strs) < 1 {
		return
	}
	isAppend := strs[0] == "-a"

	if isAppend {
		strs = strs[1:]
	}
	if len(strs) < 3 || !slices.Contains(strs, ">>") {
		for _, s := range strs {
			fmt.Println(s)
		}

		return
	}

	path := strs[len(strs)-1]
	var strRes string

	// если в аргументах есть кавычки, записываем всё что внутри них
	if slicesHasQuot(strs...) {
		strRes = strings.Join(strs[:len(strs)-2], " ")
	}

	strRes = strings.Trim(strRes, "\"")

	// если в аргументах есть кавычки пишем всчё что внутри низ в файл
	if strRes != "" {
		err := std.writeToFile(path, strRes, isAppend)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	// иначе пишем все строки с новой строки в файл
	for i, s := range strs {
		if i < len(strs)-2 {
			err := std.writeToFile(path, s+"\n", isAppend)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// writeToFile - метод для записи в файл
func (std *Stdout) writeToFile(path, str string, isAppend bool) error {
	flags := os.O_WRONLY | os.O_TRUNC | os.O_CREATE

	if isAppend {
		flags = os.O_WRONLY | os.O_APPEND | os.O_CREATE
		str += "\n"
	}

	file, errF := os.OpenFile(
		path,
		flags,
		0666,
	)

	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Fatalf("write to file: can not close file: %s", err.Error())
		}
	}(file)

	if errF != nil {
		return fmt.Errorf("write to file: %s", errF.Error())
	}

	_, errW := file.Write([]byte(str))

	if errW != nil {
		return fmt.Errorf("write to file: %s", errW.Error())
	}

	return nil
}

// kill - метод убивающий процесс по его pid
func (p *Processes) kill(pidStr string) error {
	pid, errConv := strconv.Atoi(pidStr)

	if errConv != nil {
		return fmt.Errorf("kill: illegal pid: %s", pidStr)
	}

	prcss, errPrc := os.FindProcess(pid)

	if errPrc != nil {
		return fmt.Errorf("kill: kill %d failed: no such process", pid)
	}

	errKill := prcss.Kill()

	if errKill != nil {
		return fmt.Errorf("kill: kill %d failed: %s", pid, errKill.Error())
	}

	return nil
}

// ps - выводим список запущеных процессов
func (p *Processes) ps(args []string) error {
	prcs, err := exec.Command("ps", args...).Output()

	if err != nil {
		return fmt.Errorf("process: %s", err.Error())
	}

	fmt.Print(strings.Trim(string(prcs), " "))
	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	_ = cancel
	fmt.Println("Welcome to my shell!")
	cmdLineCh := make(chan string)

	path, errPath := os.Getwd()

	if errPath != nil {
		fmt.Println("shell can not start")
	}

	dir := NewDirectories(path)
	std := NewStdout()
	prc := NewProcesses()
	shell := NewShell(dir, std, prc)

	reader := bufio.NewReader(os.Stdin)

	go func() {
		<-ctx.Done()
		close(cmdLineCh)
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				line, _, _ := reader.ReadLine()
				cmdLineCh <- string(line)
			}
		}
	}()

	go func() {
		for {
			for {
				select {
				case <-ctx.Done():
					return
				case val := <-cmdLineCh:
					if val == "quit" {
						cancel()
						return
					}

					if len(val) != 0 {
						shell.Execute(val)
					}
				}
			}
		}
	}()

	<-ctx.Done()
}

type DirWorker interface {
	changeDir(dir string)
	pwd()
}

type StdWorker interface {
	echo(str ...string)
}

type PrcWorker interface {
	kill(pidStr string) error
	ps(args []string) error
}

func slicesHasQuot(strs ...string) bool {
	if slices.Contains(strs, "\"") {
		return true
	}

	for _, str := range strs {
		if strings.Contains(str, "\"") {
			return true
		}
	}

	return false
}
