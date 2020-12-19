package demo

import (
	"app/internal/demo/router/v1"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// initWeb 初始化web引擎
func initDemoWeb() (*gin.Engine, error) {
	gin.SetMode(viper.GetString("runmode"))
	app := gin.New()
	// 注册/api路由
	err := v1.RegisterRouter(app)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func initDebugWeb() (*gin.Engine, error) {
	app := gin.New()
	pprof.Register(app)
	return app, nil
}
