package service

import (
	pb "app/api/blog/v1"
	"app/internal/biz"
	"context"
)

type PostService struct {
	pb.UnimplementedPostServer
	b *biz.PostBiz
}

func NewPostService(b *biz.PostBiz) *PostService {
	return &PostService{
		pb.UnimplementedPostServer{},
		b,
	}
}

//DTO -> DO, DO -> DTO
func (p PostService) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.PostReply, error) {
	reply := &pb.PostReply{}
	post, err := p.b.GetPost(req.GetId())
	if err != nil {
		return reply, err
	}
	reply.Id = post.Id
	reply.Title = post.Title
	return reply, nil
}
