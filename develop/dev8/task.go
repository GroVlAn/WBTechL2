package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

type Directories struct {
	curDir   strings.Builder
	backup   string
	homePath string
}

type Stdout struct{}

type Processes struct{}

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

func (sh *Shell) Execute(cmdArgs string) {

	cmds := strings.Split(cmdArgs, "|")

	for _, cmd := range cmds {
		cmd = strings.Trim(cmd, " ")
		var args []string
		lineParts := strings.Split(cmd, " ")
		cmd := lineParts[0]
		if len(lineParts) > 1 {
			args = lineParts[1:]
		}
		sh.executeCommand(cmd, args)
	}

}

func (sh *Shell) executeCommand(cmd string, args []string) {
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
		for _, arg := range args {
			sh.std.echo(arg)
		}
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

func (sh *Shell) exec(cmd string, args []string) error {
	command, err := exec.Command(cmd, args...).Output()

	if err != nil {
		return fmt.Errorf("exec: %s\n", err.Error())
	}

	fmt.Print(string(command))
	return nil
}

func (d *Directories) changeDir(dir string) {
	if len(dir) == 0 {
		d.curDir.Reset()
		d.curDir.WriteString(d.homePath)
		return
	}
	parts := strings.Split(dir, "/")
	isNext := false
	isResetNotDots := false

	d.backup = d.curDir.String()

	for _, part := range parts {
		if part == "" {
			continue
		}
		curDir := d.curDir.String()
		if part == ".." {
			lstSlh := strings.LastIndex(curDir, "/")

			if lstSlh == -1 {
				continue
			}
			d.curDir.Reset()
			d.curDir.WriteString(curDir[:lstSlh])
			continue
		}
		if part == "." {
			isNext = true
			continue
		} else if !isResetNotDots && !isNext {
			d.curDir.Reset()
			isResetNotDots = true
		}
		d.curDir.WriteString("/")
		d.curDir.WriteString(part)
	}

	if !d.exist() {
		fmt.Printf("cd: path %s does not exist\n", d.curDir.String())
		d.curDir.Reset()
		d.curDir.WriteString(d.backup)
	}
}

func (d *Directories) exist() bool {
	_, err := os.Stat(d.curDir.String())

	if err != nil {
		return false
	}

	return true
}

func (d *Directories) pwd() {
	fmt.Println(d.curDir.String())
}

func (std *Stdout) echo(str string) {
	fmt.Println(str)
}

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
	echo(str string)
}

type PrcWorker interface {
	kill(pidStr string) error
	ps(args []string) error
}
