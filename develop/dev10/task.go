package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	fTimeout := flag.Int("t", 10, "Connection end time in seconds")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("error: empty ip and port")
	}
	// Создаем TCP-соединение
	addr := buildAddress(args)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("can not close connection", err.Error())
		}
	}(conn)

	// Выход по Ctrl + С
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	go func() {
		<-time.After(time.Duration(*fTimeout) * time.Second)
		gracefulShutdown <- os.Interrupt
	}()
	// Общение клиента с сервером
	go copyTo(gracefulShutdown, os.Stdout, conn) // читаем из сокета
	go copyTo(gracefulShutdown, conn, os.Stdin)  // пишем в сокет

	<-gracefulShutdown
	log.Println("connection was closed")
}

func copyTo(gracefulShutdown chan os.Signal, dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Println(err)
		gracefulShutdown <- os.Interrupt
	}
}

func buildAddress(args []string) string {
	var b strings.Builder

	b.WriteString(args[0])
	b.WriteString(":")
	b.WriteString(args[1])

	return b.String()
}
