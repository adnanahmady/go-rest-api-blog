package application

import (
	"context"
	"errors"

	"github.com/adnanahmady/go-rest-api-blog/internal/domain"
	"github.com/adnanahmady/go-rest-api-blog/pkg/errs"
)

type ShowPostUseCase interface {
	Show(ctx context.Context, id uint) (*PostDTO, error)
}

var _ ShowPostUseCase = (*ShowPostUseCaseImpl)(nil)

type ShowPostUseCaseImpl struct {
	repo domain.PostRepository
}

func NewShowPostUseCase(repo domain.PostRepository) *ShowPostUseCaseImpl {
	return &ShowPostUseCaseImpl{repo: repo}
}

func (uc *ShowPostUseCaseImpl) Show(
	ctx context.Context,
	id uint,
) (*PostDTO, error) {
	post, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrPostNotFound) {
			return nil, errs.NewNotFoundError("post not found")
		}
		return nil, err
	}
	return newPostDTO(post), nil
}
