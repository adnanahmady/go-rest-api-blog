package application

import (
	"context"
	"time"

	"github.com/adnanahmady/go-rest-api-blog/internal/domain"
)

type CreatePostUseCase interface {
	Create(ctx context.Context, title, content string) (*PostDTO, error)
}

var _ CreatePostUseCase = (*CreatePostUseCaseImpl)(nil)

type CreatePostUseCaseImpl struct {
	repo domain.PostRepository
}

func NewCreatePostUseCase(repo domain.PostRepository) *CreatePostUseCaseImpl {
	return &CreatePostUseCaseImpl{repo: repo}
}

func (uc *CreatePostUseCaseImpl) Create(
	ctx context.Context,
	title, content string,
) (*PostDTO, error) {
	post := &domain.Post{
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := uc.repo.Create(ctx, post); err != nil {
		return nil, err
	}
	return newPostDTO(post), nil
}