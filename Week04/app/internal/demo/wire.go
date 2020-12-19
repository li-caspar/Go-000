// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package demo

import "github.com/google/wire"

func InitializeApp(filename string) *App {
	wire.Build(NewConfig, NewLogger, NewDB, NewEngine, NewApp)
	return &App{}
}
