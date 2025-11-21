package infra

import (
	"context"
	"database/sql"
	"errors"

	"github.com/adnanahmady/go-rest-api-blog/internal/domain"
	"github.com/adnanahmady/go-rest-api-blog/pkg/database"
	"github.com/adnanahmady/go-rest-api-blog/pkg/request"
)

var _ domain.PostRepository = (*SqlitePostRepository)(nil)

type SqlitePostRepository struct {
	dbm database.DatabaseManager
}

func NewSqlitePostRepository(
	dbm database.DatabaseManager,
) *SqlitePostRepository {
	return &SqlitePostRepository{dbm: dbm}
}

func (r *SqlitePostRepository) Create(
	ctx context.Context,
	post *domain.Post,
) error {
	lgr := request.GetLogger(ctx)

	var dbName string
	err := r.dbm.GetClient().GetContext(ctx, &dbName, `
	PRAGMA database_list; select name from pragma_database_list limit 1`)
	if err != nil {
		lgr.Error("failed to get database name", err)
		return err
	}
	lgr.Info("database name", "database", dbName)
	query := `INSERT INTO posts (title, content, created_at, updated_at)
	VALUES (:title, :content, :created_at, :updated_at)`

	_, err = r.dbm.GetClient().NamedExecContext(ctx, query, post)
	if err != nil {
		lgr.Error("failed to create post", err)
		return err
	}
	lgr.Info("post created successfully")
	return nil
}

func (r *SqlitePostRepository) GetByID(
	ctx context.Context,
	id uint,
) (*domain.Post, error) {
	var post domain.Post
	query := `SELECT * FROM posts WHERE id = ?`
	err := r.dbm.GetClient().GetContext(ctx, &post, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPostNotFound
		}
		return nil, err
	}
	return &post, nil
}

func (r *SqlitePostRepository) Update(
	ctx context.Context,
	post *domain.Post,
) error {
	query := `UPDATE posts SET title = :title, content = :content,
		updated_at = :updated_at WHERE id = :id`
	_, err := r.dbm.GetClient().NamedExecContext(ctx, query, post)
	if err != nil {
		return err
	}
	return err
}

func (r *SqlitePostRepository) Delete(
	ctx context.Context, id uint,
) error {
	query := `DELETE FROM posts WHERE id = ?`
	_, err := r.dbm.GetClient().ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return err
}

func (r *SqlitePostRepository) GetPaginated(
	ctx context.Context,
	page, perPage int,
) ([]*domain.Post, int, error) {
	var posts []*domain.Post
	query := `SELECT * FROM posts LIMIT ? OFFSET ?`
	err := r.dbm.GetClient().Select(&posts, query, perPage, (page-1)*perPage)
	if err != nil {
		return nil, 0, err
	}

	var total int
	query = `SELECT COUNT(*) as total FROM posts LIMIT ? OFFSET ?`
	err = r.dbm.GetClient().GetContext(ctx, &total, query, perPage, (page-1)*perPage)
	if err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}
