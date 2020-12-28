// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"app/internal/biz"
	"app/internal/config"
	"app/internal/data"
	"app/internal/service"
	"github.com/google/wire"
)

func InitializePostService(cfg *config.Config) (*service.PostService, error) {
	wire.Build(
		data.NewDB,
		data.NewPostData,
		wire.Bind(new(biz.PostRepo), new(*data.PostData)),
		biz.NewPostBiz,
		service.NewPostService)
	return &service.PostService{}, nil
}
