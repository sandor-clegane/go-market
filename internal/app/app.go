package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sandor-clegane/go-market/internal/config"
	"github.com/sandor-clegane/go-market/internal/handlers/orderhandler"
	"github.com/sandor-clegane/go-market/internal/handlers/userhandler"
	withdrawHandler2 "github.com/sandor-clegane/go-market/internal/handlers/withdrawhandler"
	"github.com/sandor-clegane/go-market/internal/router"
	cookieService2 "github.com/sandor-clegane/go-market/internal/service/cookieservice"
	orderService2 "github.com/sandor-clegane/go-market/internal/service/orderservice"
	userService2 "github.com/sandor-clegane/go-market/internal/service/userservice"
	withdrawService2 "github.com/sandor-clegane/go-market/internal/service/withdrawservice"
	"github.com/sandor-clegane/go-market/internal/storage"
	"github.com/sandor-clegane/go-market/internal/storage/orderstorage"
	"github.com/sandor-clegane/go-market/internal/storage/userstorage"
	"github.com/sandor-clegane/go-market/internal/storage/withdrawstorage"
)

const (
	readTimeout  = 10 * time.Second
	writeTimeout = 10 * time.Second
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
	userStg := userstorage.New(db)
	orderStg, err := orderstorage.New(db, cfg.AccrualSystemAddress)
	if err != nil {
		return nil, err
	}
	withdrawStg := withdrawstorage.New(db)
	//create services
	userService := userService2.New(userStg)
	orderService := orderService2.New(orderStg)
	withdrawService := withdrawService2.New(withdrawStg, orderStg)
	cookieService, err := cookieService2.New(cfg.Key)
	if err != nil {
		return nil, err
	}
	//create handlers
	urlHandler := userhandler.New(userService, cookieService)
	orderHandler := orderhandler.New(orderService, cookieService)
	withdrawHandler := withdrawHandler2.New(withdrawService, cookieService)
	//create router
	urlRouter := router.NewRouter(urlHandler, orderHandler, withdrawHandler)

	server := &http.Server{
		Addr:         cfg.ServerEndpoint,
		Handler:      urlRouter,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	defer listenForStorageCloseSignal(orderStg)
	return &App{
		HTTPServer: server,
	}, nil
}

func listenForStorageCloseSignal(repository orderstorage.OrderStorage) {
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
