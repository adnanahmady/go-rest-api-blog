package application

import (
	"math"
	"time"

	"github.com/adnanahmady/go-rest-api-blog/internal/domain"
)

type PostDTO struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newPostDTO(post *domain.Post) *PostDTO {
	return &PostDTO{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}

type PaginationDTO struct {
	Data       []*PostDTO `json:"-"`
	Total      int        `json:"total"`
	Page       int        `json:"page"`
	PerPage    int        `json:"per_page"`
	TotalPages int        `json:"total_pages"`
}

func newPaginationDTO(
	posts []*PostDTO,
	total, page, perPage int,
) *PaginationDTO {
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &PaginationDTO{
		Data:       posts,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}
}
