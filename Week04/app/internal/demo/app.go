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
}

var DemoApp *App

func NewApp(cfg Config, logger *Logger, db *DataBase) *App {
	if DemoApp != nil {
		return DemoApp
	}
	DemoApp = &App{
		config:   cfg,
		DataBase: db,
		Logger:   logger,
	}
	return DemoApp
}

func (app App) Start() error {
	g, ctx := errgroup.WithContext(context.Background())
	//启动http服务
	g.Go(func() error {
		demoweb, err := initDemoWeb()
		if err != nil {
			return err
		}
		addr := fmt.Sprintf("%s:%d", viper.GetString("http.host"), viper.GetInt("http.port"))
		return app.StartHTTPServer(ctx, addr, demoweb)
	})
	//启动debug服务
	g.Go(func() error {
		demoweb, err := initDebugWeb()
		if err != nil {
			return err
		}
		addr := fmt.Sprintf("%s:%d", viper.GetString("debug.host"), viper.GetInt("debug.port"))
		return app.StartHTTPServer(ctx, addr, demoweb)
	})
	//启动redis缓存服务
	/*g.Go(func() error {

	})*/

	sc := make(chan os.Signal)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//根据信号量来退出
	g.Go(func() error {
		return serverSignal(ctx, sc)
	})
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

//启动指定的服务
func (app *App) StartHTTPServer(ctx context.Context, addr string, handler *gin.Engine) error {
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		select {
		case <-ctx.Done():
			ctxTimeOut, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(viper.GetInt("http.shutdown_timeout")))
			defer cancel()
			if err := srv.Shutdown(ctxTimeOut); err != nil {
				app.Logger.Errorf(err.Error())
			}
		}
	}()
	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

//监听信号量
func serverSignal(ctx context.Context, sc chan os.Signal) error {
	select {
	case sig := <-sc:
		return fmt.Errorf("signal exit, %s", sig.String())
	case <-ctx.Done():
		fmt.Println("context done, signa return")
		return nil
	}
}
