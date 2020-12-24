package data

import "app/internal/domain"

type PostData struct {
	//db
	//redis
}

func NewPostData() *PostData {
	return &PostData{}
}

func (p PostData) GetPost(id int64) (domain.Post, error) {
	return domain.Post{
		Id:    id,
		Title: "test",
	}, nil

}
