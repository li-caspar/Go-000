package demo

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//当前应用名称
var AppName string = "demo"

type App struct {
	config   Config
	DataBase *DataBase
	Logger   *Logger
	Engine   *Engine
}

var DemoApp *App

func NewApp(cfg Config, logger *Logger, db *DataBase, engine *Engine) *App {
	if DemoApp != nil {
		return DemoApp
	}
	DemoApp = &App{
		config:   cfg,
		DataBase: db,
		Logger:   logger,
		Engine:   engine,
	}
	return DemoApp
}

func (app *App) StartHTTPServer(ctx context.Context, handler *gin.Engine) error {
	addr := fmt.Sprintf("%s:%d", viper.GetString("http.host"), viper.GetInt("http.port"))
	sx := Serverx{}
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	cancel := func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(viper.GetInt("http.shutdown_timeout")))
		defer cancel()
		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			app.Logger.Errorf(err.Error())
		}
	}
	sx.server = srv
	sx.cancel = cancel
	app.Engine.AddServerxs(sx)
	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (app App) Start() error {
	//var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		web := InitWeb()
		return app.StartHTTPServer(ctx, web)
	})

	sc := make(chan os.Signal)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//根据信号量来退出
	g.Go(func() error {
		return serverSignal(ctx, sc)
	})
	return g.Wait()
}

func serverSignal(ctx context.Context, sc chan os.Signal) error {
	select {
	case sig := <-sc:
		return fmt.Errorf("signal exit, %s", sig.String())
	case <-ctx.Done():
		fmt.Println("context done, signa return")
		return nil
	}
}
