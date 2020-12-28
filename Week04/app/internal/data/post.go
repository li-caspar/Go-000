package data

import (
	"app/internal/domain"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type PostData struct {
	db *gorm.DB
	//redis
}

func NewPostData(db *gorm.DB) *PostData {
	return &PostData{
		db: db,
	}
}

func (p PostData) GetPost(id int64) (domain.Post, error) {
	post := domain.Post{}
	if err := p.db.First(&post).Error; err != nil {
		return post, errors.Wrap(err, "db first post error")
	}
	return post, nil
}
