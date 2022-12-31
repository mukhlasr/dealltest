package posts

import (
	"context"
	"time"
)

type StoredPost struct {
	ID        int64
	Title     string
	Content   string
	Timestamp time.Time
}

func (s StoredPost) isValid() bool {
	if len(s.Title) < 3 {
		return false
	}

	return true
}

type PostAccessor interface {
	GetAll(ctx context.Context) ([]StoredPost, error)
	GetByID(ctx context.Context, id int64) (StoredPost, error)
}

type PostMutator interface {
	Update(ctx context.Context, p StoredPost) error
	DeletePostByID(ctx context.Context, id int64) error
}

type PostCreator interface {
	Create(ctx context.Context, p StoredPost) (id int64, err error)
}

type PostStorage interface {
	PostAccessor
	PostCreator
	PostMutator
}
