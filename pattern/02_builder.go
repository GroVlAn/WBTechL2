package pattern

import "fmt"

/*
Преймущества:
	Позволяет создавать куски программы пошагово
	Позволяет использовать использовать один и тот же кусок кода для создания различных кусков программы
	Изолирует сборку программы от бизнес логики
Недостатки:
	Усложняет код из за добавления дополнительных структур
	Тот кто обращается к строителю будет привязан к интерфесу стоителя, так как директор может не иметь метода
	получения данных
*/

type Config struct{}

type Logger struct{}

type Repository struct{}

type Service struct{}

type Handler struct{}

type Server struct{}

type ApplicationBuilder struct {
	config  Config
	logger  Logger
	repo    Repository
	service Service
	handler Handler
	server  Server
}

type ApplicationDirector struct{}

func (cfg *Config) InitConfig() {
	fmt.Println("initializing config...")
}

func (cfg *Config) Config() {
	fmt.Println("get current config")
}

func (lg *Logger) InitLogger() {
	fmt.Println("initializing logger...")
}

func (lg *Logger) Logger() {
	fmt.Println("get current logger")
}

func (rep *Repository) InitRepository() {
	fmt.Println("initializing repository...")
}

func (rep *Repository) Repository() {
	fmt.Println("get current repository")
}

func (srv *Service) InitService() {
	fmt.Println("initializing service...")
}

func (srv *Service) Service() {
	fmt.Println("get current service")
}

func (hnd *Handler) InitHandler() {
	fmt.Println("initializing handler...")
}

func (hnd *Handler) Handler() {
	fmt.Println("get handler")
}

func (svr *Server) InitServer() {
	fmt.Println("initializing server...")
}

func (svr *Server) Server() {
	fmt.Println("get server")
}

func (ab *ApplicationBuilder) SetConfig() {
	ab.config.InitConfig()
	ab.config.Config()
}

func (ab *ApplicationBuilder) SetLogger() {
	ab.logger.InitLogger()
	ab.logger.Logger()
}

func (ab *ApplicationBuilder) SetRepository() {
	ab.repo.InitRepository()
	ab.repo.Repository()
}

func (ab *ApplicationBuilder) SetService() {
	ab.service.InitService()
	ab.service.Service()
}

func (ab *ApplicationBuilder) SetHandler() {
	ab.handler.InitHandler()
	ab.handler.Handler()
}

func (ab *ApplicationBuilder) SetServer() {
	ab.server.InitServer()
	ab.server.Server()
	ab.server.Server()
}

func (ad *ApplicationDirector) BuildApp(builder Builder) {
	builder.SetConfig()
	builder.SetLogger()
	builder.SetRepository()
	builder.SetService()
	builder.SetHandler()
	builder.SetServer()
}

type Builder interface {
	SetConfig()
	SetLogger()
	SetRepository()
	SetService()
	SetHandler()
	SetServer()
}
