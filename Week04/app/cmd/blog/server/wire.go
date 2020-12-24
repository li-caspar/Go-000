// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"app/internal/biz"
	"app/internal/data"
	"app/internal/service"
	"github.com/google/wire"
)

func InitializePostService() *service.PostService {
	wire.Build(
		data.NewPostData,
		wire.Bind(new(biz.PostRepo), new(*data.PostData)),
		biz.NewPostBiz,
		service.NewPostService)
	return &service.PostService{}
}
