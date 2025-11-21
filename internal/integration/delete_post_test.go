package integration

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/adnanahmady/go-rest-api-blog/internal/domain"
	"github.com/adnanahmady/go-rest-api-blog/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestDeletePost(t *testing.T) {
	server, err := test.Setup()
	assert.NoError(t, err)
	defer server.Close()
	postFactory := NewPostFactory(server.App.Database)

	name := "given post when requested to delete then it should delete the post"
	t.Run(name, func(t *testing.T) {
		// Arrange
		post, err := postFactory.CreatePost()
		assert.NoError(t, err)
		url := fmt.Sprintf("/v1/posts/%d", post.ID)

		// Act
		server.Delete(t, url, nil)

		// Assert
		err = server.App.Database.GetClient().Get(
			&domain.Post{},
			"SELECT * FROM posts WHERE id = ?",
			post.ID,
		)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	name = "given post when requested but not found then return 404"
	t.Run(name, func(t *testing.T) {
		// Arrange
		url := "/v1/posts/99999"

		// Act
		rec, _ := server.Delete(t, url, nil)

		// Assert
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	name = "given post when requested then return 204"
	t.Run(name, func(t *testing.T) {
		// Arrange
		post, err := postFactory.CreatePost()
		assert.NoError(t, err)
		url := fmt.Sprintf("/v1/posts/%d", post.ID)

		// Act
		rec, _ := server.Delete(t, url, nil)

		// Assert
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}
