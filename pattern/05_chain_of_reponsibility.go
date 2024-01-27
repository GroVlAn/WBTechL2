package pattern

import (
	"fmt"
	"net/http"
)

/*
Преймущества:
	- Уменьшает зависимость между клиентом и обработчиками.
	- Реализует принцип единственной обязанности.
	- Реализует принцип открытости/закрытости.
Недостатки:
	- Запрос может остаться никем не обработанным.
*/

// AuthHandler - проверяет аунтентификацию пользователея
type AuthHandler struct {
	nextHandler HttpHandler
}

// LoggingHandler - обработчик логов
type LoggingHandler struct {
	nextHandler HttpHandler
}

// MainHandler - основной обработчик запросов
type MainHandler struct {
	nextHandler HttpHandler
}

// HandlerRequest - обрабатываем запрос на аунтентификацию и если если следующий обработчик передаём запрос ему
func (ah *AuthHandler) HandlerRequest(req *http.Request) bool {
	fmt.Println("AuthHandler: checking authentication ...")

	if ah.nextHandler != nil {
		return ah.nextHandler.HandlerRequest(req)
	}

	return true
}

// SetNext - добавляем следующий обработчик
func (ah *AuthHandler) SetNext(handler HttpHandler) {
	ah.nextHandler = handler
}

// HandlerRequest - выполняет логирование запроса и если есть следующий обработчик передаёт запрос ему
func (lh *LoggingHandler) HandlerRequest(req *http.Request) bool {
	fmt.Println("LoggingHandler: logging request ...")

	if lh.nextHandler != nil {
		return lh.nextHandler.HandlerRequest(req)
	}

	return true
}

func (lh *LoggingHandler) SetNext(handler HttpHandler) {
	lh.nextHandler = handler
}

// HandlerRequest - обрабатываем главный запрос
func (mh *MainHandler) HandlerRequest(req *http.Request) bool {
	fmt.Println("MainHandler: handling main request ...")

	return true
}

func (mh *MainHandler) SetNext(handler HttpHandler) {
	mh.nextHandler = handler
}

// HttpHandler - интерфейс для обработчиков запросов
type HttpHandler interface {
	HandlerRequest(req *http.Request) bool
	SetNext(handler HttpHandler)
}
