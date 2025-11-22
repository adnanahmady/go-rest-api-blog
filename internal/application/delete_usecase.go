package application

import (
	"context"
	"errors"

	"github.com/adnanahmady/go-rest-api-blog/internal/domain"
	"github.com/adnanahmady/go-rest-api-blog/pkg/errs"
)

type DeletePostUseCase interface {
	Delete(ctx context.Context, id uint) error
}

var _ DeletePostUseCase = (*DeletePostUseCaseImpl)(nil)

type DeletePostUseCaseImpl struct {
	repo domain.PostRepository
}

func NewDeletePostUseCase(repo domain.PostRepository) *DeletePostUseCaseImpl {
	return &DeletePostUseCaseImpl{repo: repo}
}

func (uc *DeletePostUseCaseImpl) Delete(
	ctx context.Context,
	id uint,
) error {
	if _, err := uc.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, domain.ErrPostNotFound) {
			return errs.NewNotFoundError("post not found")
		}
		return err
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
