package demo

import (
	v1 "app/internal/demo/router/v1"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

type Serverx struct {
	server *http.Server
	cancel func()
}

type Engine struct {
	Serverxs []Serverx
}

func NewEngine() *Engine {
	engine := &Engine{
		Serverxs: make([]Serverx, 0),
	}
	return engine
}

func (d *Engine) Stop() {
	if d == nil {
		return
	}
	if len(d.Serverxs) > 0 {
		for _, serverx := range d.Serverxs {
			if serverx.cancel != nil {
				serverx.cancel()
			}
		}
	}
}

func (d *Engine) AddServerxs(serverx Serverx) {
	d.Serverxs = append(d.Serverxs, serverx)
}

// InitWeb 初始化web引擎
func InitWeb() *gin.Engine {
	gin.SetMode(viper.GetString("runmode"))
	app := gin.New()
	// 注册/api路由
	err := v1.RegisterRouter(app)
	if err != nil {
		panic(err)
	}
	return app
}
