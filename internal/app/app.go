package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Gophermarket/go-musthave-diploma-tpl/internal/config"
	"Gophermarket/go-musthave-diploma-tpl/internal/handlers/ordersHandler"
	"Gophermarket/go-musthave-diploma-tpl/internal/handlers/userHandler"
	withdrawHandler2 "Gophermarket/go-musthave-diploma-tpl/internal/handlers/withdrawHandler"
	"Gophermarket/go-musthave-diploma-tpl/internal/router"
	cookieService2 "Gophermarket/go-musthave-diploma-tpl/internal/service/cookieService"
	orderService2 "Gophermarket/go-musthave-diploma-tpl/internal/service/orderService"
	userService2 "Gophermarket/go-musthave-diploma-tpl/internal/service/userService"
	withdrawService2 "Gophermarket/go-musthave-diploma-tpl/internal/service/withdrawService"
	"Gophermarket/go-musthave-diploma-tpl/internal/storage"
	"Gophermarket/go-musthave-diploma-tpl/internal/storage/orderStorage"
	"Gophermarket/go-musthave-diploma-tpl/internal/storage/userStorage"
	"Gophermarket/go-musthave-diploma-tpl/internal/storage/withdrawStorage"
)

const (
	readTimeout    = 10 * time.Second
	writeTimeout   = 10 * time.Second
	maxHeaderBytes = 1 << 20
)

type App struct {
	HTTPServer *http.Server
}

func New(cfg config.Config) (*App, error) {
	log.Println("creating router")
	db, err := storage.ConnectAndInitDB(cfg.DatabaseAddress)
	if err != nil {
		return nil, err
	}
	//create storages
	userRepository := userStorage.New(db)
	orderRepository, err := orderStorage.New(db, cfg.AccrualSystemAddress)
	if err != nil {
		return nil, err
	}
	withdrawRepository := withdrawStorage.New(db)
	//create services
	userService := userService2.New(userRepository)
	orderService := orderService2.New(orderRepository)
	withdrawService := withdrawService2.New(withdrawRepository, orderRepository)
	cookieService, err := cookieService2.New(cfg.Key)
	if err != nil {
		return nil, err
	}
	//create handlers
	urlHandler := userHandler.New(userService, cookieService)
	orderHandler := ordersHandler.New(orderService, cookieService)
	withdrawHandler := withdrawHandler2.New(withdrawService, cookieService)
	//create router
	urlRouter := router.NewRouter(urlHandler, orderHandler, withdrawHandler)

	server := &http.Server{
		Addr:           cfg.ServerEndpoint,
		Handler:        urlRouter,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}
	defer listenForStorageCloseSignal(server, orderRepository)
	return &App{
		HTTPServer: server,
	}, nil
}

func listenForStorageCloseSignal(server *http.Server, repository orderStorage.OrderStorage) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		repository.StopSchedulerAndWorkerPool()
	}()
}

func (app *App) Start() error {
	log.Println("start server")
	return app.HTTPServer.ListenAndServe()
}

func (app *App) Shutdown(ctx context.Context) error {
	return app.HTTPServer.Shutdown(ctx)
}
