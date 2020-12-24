package biz

import (
	"app/internal/domain"
	"fmt"
)

type PostRepo interface {
	GetPost(int64) (domain.Post, error)
}

func NewPostBiz(repo PostRepo) *PostBiz {
	return &PostBiz{
		repo,
	}
}

type PostBiz struct {
	d PostRepo
}

func (p PostBiz) GetPost(id int64) (domain.Post, error) {
	var post domain.Post
	if id < 1 {
		return post, fmt.Errorf("[biz] GetPost GetPost error:%d", id)
	}
	post, err := p.d.GetPost(id)
	return post, err
}
