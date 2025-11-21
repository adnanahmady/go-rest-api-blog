package application

import (
	"context"

	"github.com/adnanahmady/go-rest-api-blog/internal/domain"
)

type ListPostsUseCase interface {
	List(ctx context.Context, page, perPage int) (*PaginationDTO, error)
}

var _ ListPostsUseCase = (*ListPostsUseCaseImpl)(nil)

type ListPostsUseCaseImpl struct {
	repo domain.PostRepository
}

func NewListPostsUseCase(repo domain.PostRepository) *ListPostsUseCaseImpl {
	return &ListPostsUseCaseImpl{repo: repo}
}

func (uc *ListPostsUseCaseImpl) List(
	ctx context.Context,
	page, perPage int,
) (*PaginationDTO, error) {
	posts, total, err := uc.repo.GetPaginated(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	collection := make([]*PostDTO, len(posts))
	for i, post := range posts {
		collection[i] = newPostDTO(post)
	}

	return newPaginationDTO(collection, total, page, perPage), nil
}
