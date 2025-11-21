package application

import (
	"context"
	"time"

	"github.com/adnanahmady/go-rest-api-blog/internal/domain"
)

type UpdatePostUseCase interface {
	Update(ctx context.Context, id uint, title, content string) (*PostDTO, error)
}

var _ UpdatePostUseCase = (*UpdatePostUseCaseImpl)(nil)

type UpdatePostUseCaseImpl struct {
	repo domain.PostRepository
}

func NewUpdatePostUseCase(repo domain.PostRepository) *UpdatePostUseCaseImpl {
	return &UpdatePostUseCaseImpl{repo: repo}
}

func (uc *UpdatePostUseCaseImpl) Update(
	ctx context.Context,
	id uint,
	title, content string,
) (*PostDTO, error) {
	post := &domain.Post{
		ID:        id,
		Title:     title,
		Content:   content,
		UpdatedAt: time.Now(),
	}
	if err := uc.repo.Update(ctx, post); err != nil {
		return nil, err
	}
	return newPostDTO(post), nil
}
