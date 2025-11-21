package domain

import "context"

type PostRepository interface {
	Create(ctx context.Context, post *Post) error
	GetByID(ctx context.Context, id uint) (*Post, error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id uint) error
	GetPaginated(ctx context.Context, page, perPage int) ([]*Post, int, error)
}
