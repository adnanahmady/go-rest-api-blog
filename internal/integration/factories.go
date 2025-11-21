package integration

import (
	"context"
	"time"

	"github.com/adnanahmady/go-rest-api-blog/internal/domain"
	"github.com/adnanahmady/go-rest-api-blog/pkg/database"
	"github.com/brianvoe/gofakeit/v6"
)

type PostFactory struct {
	dbm database.DatabaseManager
}

func NewPostFactory(dbm database.DatabaseManager) *PostFactory {
	return &PostFactory{dbm: dbm}
}

func (f *PostFactory) CreatePost(fn ...func(*domain.Post)) (*domain.Post, error) {
	post := &domain.Post{
		Title:     gofakeit.Sentence(10),
		Content:   gofakeit.Paragraph(10, 10, 10, " "),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	for _, f := range fn {
		f(post)
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	query := "INSERT INTO posts (title, content, created_at, updated_at) " +
		"VALUES (:title, :content, :created_at, :updated_at)"
	res, err := f.dbm.GetClient().NamedExecContext(ctx, query, post)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	post.ID = uint(id)
	return post, nil
}
