package service

import (
	pb "app/api/blog/v1"
	"context"
)

type PostService struct {
	pb.UnimplementedPostServer
}

//DTO -> DO, DO -> DTO
func (p PostService) GetPost(context.Context, *pb.GetPostRequest) (*pb.PostReply, error) {

	return &pb.PostReply{
		Id:    1,
		Title: "test",
	}, nil
}
