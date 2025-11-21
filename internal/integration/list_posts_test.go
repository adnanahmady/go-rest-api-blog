package integration

import (
	"net/http"
	"testing"

	"github.com/adnanahmady/go-rest-api-blog/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestListPosts(t *testing.T) {
	server, err := test.Setup()
	assert.NoError(t, err)
	defer server.Close()
	postFactory := NewPostFactory(server.App.Database)

	name := "given posts when requested then it should list them"
	t.Run(name, func(t *testing.T) {
		// Arrange
		server.App.Database.GetClient().Exec("DELETE FROM posts")
		defer func() {
			server.App.Database.GetClient().Exec("DELETE FROM posts")
		}()
		post, err := postFactory.CreatePost()
		assert.NoError(t, err)
		url := "/v1/posts"

		// Act
		_, body := server.Get(t, url, nil)

		// Assert
		bodyMap := body.(map[string]any)
		dataMap := bodyMap["data"].([]any)
		itemMap := dataMap[0].(map[string]any)

		assert.IsType(t, float64(0), itemMap["id"])
		assert.Equal(t, post.Title, itemMap["title"])
		assert.Equal(t, post.Content, itemMap["content"])
		assert.IsType(t, "string", itemMap["created_at"])
		assert.IsType(t, "string", itemMap["updated_at"])
		test.AssertTimeFormat(t, itemMap["created_at"].(string))
		test.AssertTimeFormat(t, itemMap["updated_at"].(string))
	})

	name = "given posts when requested then return pagination"
	t.Run(name, func(t *testing.T) {
		// Arrange
		server.App.Database.GetClient().Exec("DELETE FROM posts")
		defer func() {
			server.App.Database.GetClient().Exec("DELETE FROM posts")
		}()
		_, err = postFactory.CreatePost()
		_, err = postFactory.CreatePost()
		assert.NoError(t, err)
		url := "/v1/posts?page=1&per_page=10"

		// Act
		_, body := server.Get(t, url, nil)

		// Assert
		bodyMap := body.(map[string]any)
		metaMap := bodyMap["meta"].(map[string]any)
		paginationMap := metaMap["pagination"].(map[string]any)

		assert.IsType(t, float64(2), paginationMap["total"])
		assert.IsType(t, float64(1), paginationMap["page"])
		assert.IsType(t, float64(10), paginationMap["per_page"])
		assert.IsType(t, float64(1), paginationMap["total_pages"])
	})

	name = "given posts when there are no posts then return empty list"
	t.Run(name, func(t *testing.T) {
		// Arrange
		url := "/v1/posts"

		// Act
		_, body := server.Get(t, url, nil)

		// Assert
		bodyMap := body.(map[string]any)
		dataMap := bodyMap["data"].([]any)
		assert.Empty(t, dataMap)
	})

	name = "given posts when requested then return 200"
	t.Run(name, func(t *testing.T) {
		// Arrange
		url := "/v1/posts"

		// Act
		rec, _ := server.Get(t, url, nil)

		// Assert
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
