package main

import (
	"app/internal/demo"
	"fmt"
	"github.com/spf13/pflag"
)

var (
	cfg = pflag.StringP("config", "c", "./configs/demo.yaml", "apiserver config file path.")
)

func main() {
	pflag.Parse()
	App := demo.InitializeApp(*cfg)
	App.Logger.Println("app start")
	if err := App.Start(); err != nil {
		fmt.Printf("app failed:%v\n", err)
	}
	App.Logger.Println("app stop")
	App.DataBase.Stop()
	App.Logger.Stop()
}
