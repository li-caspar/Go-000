// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package demo

// Injectors from wire.go:

func InitializeApp(filename string) *App {
	config := NewConfig(filename)
	logger := NewLogger()
	dataBase := NewDB()
	app := NewApp(config, logger, dataBase)
	return app
}
