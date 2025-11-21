package presentation

import (
	"github.com/adnanahmady/go-rest-api-blog/internal/application"
)

type PostResource struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func newPostResource(post application.PostDTO) *PostResource {
	return &PostResource{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func newPostCollection(posts []*application.PostDTO) []*PostResource {
	collection := make([]*PostResource, len(posts))
	for i, post := range posts {
		collection[i] = newPostResource(*post)
	}
	return collection
}
