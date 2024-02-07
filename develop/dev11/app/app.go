package app

import (
	"context"
	"dev11/config"
	"dev11/server"
	"dev11/service"
	"dev11/transport/api"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

const (
	configPath = "configs"
	configName = "config"
)

type CalendarApp struct {
	conf    *config.Config
	handler Handler
	serv    Server
	evSv    api.EventServ
}

func NewApp() *CalendarApp {
	return &CalendarApp{}
}

func (ca *CalendarApp) Run(contxt context.Context) {
	if err := config.InitConfig(configPath, configName); err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}

	ctx, cancel := signal.NotifyContext(contxt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	_ = cancel

	ca.setConfig()
	ca.setHandler()
	ca.setHttpServer()
	ca.setEventService()
	//ca.handler.AddMiddleware()

	go func() {
		<-ctx.Done()
		err := ca.serv.Shutdown(ctx)

		if err != nil {
			log.Fatalf("")
			return
		}
	}()

	go func() {
		if err := ca.serv.Start(); err != nil {
			log.Fatalf("run: can not start server: %s", err.Error())
			return
		}
	}()

	<-ctx.Done()
}

func (ca *CalendarApp) setConfig() {
	ca.conf = config.NewConfig()
}

func (ca *CalendarApp) setHandler() {
	ca.handler = api.NewHTTPHandler(ca.evSv)
}

func (ca *CalendarApp) setHttpServer() {
	ca.serv = server.NewHttpServer(ca.conf, ca.handler.Handler())
}

func (ca *CalendarApp) setEventService() {
	ca.evSv = service.NewEventsService()
}

type Handler interface {
	Handler() http.Handler
}

type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
}
