package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/adnanahmady/go-rest-api-blog/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestShowPost(t *testing.T) {
	server, err := test.Setup()
	assert.NoError(t, err)
	defer server.Close()
	postFactory := NewPostFactory(server.App.Database)

	name := "given post when requested then it should show the post"
	t.Run(name, func(t *testing.T) {
		// Arrange
		post, err := postFactory.CreatePost()
		assert.NoError(t, err)
		url := fmt.Sprintf("/v1/posts/%d", post.ID)

		// Act
		_, body := server.Get(t, url, nil)

		// Assert
		bodyMap := body.(map[string]any)
		dataMap := bodyMap["data"].(map[string]any)

		assert.IsType(t, float64(0), dataMap["id"])
		assert.Equal(t, post.Title, dataMap["title"])
		assert.Equal(t, post.Content, dataMap["content"])
		assert.IsType(t, "string", dataMap["created_at"])
		assert.IsType(t, "string", dataMap["updated_at"])
		test.AssertTimeFormat(t, dataMap["created_at"].(string))
		test.AssertTimeFormat(t, dataMap["updated_at"].(string))
	})

	name = "given post when requested but not found then return 404"
	t.Run(name, func(t *testing.T) {
		// Arrange
		url := "/v1/posts/99999"

		// Act
		rec, _ := server.Get(t, url, nil)

		// Assert
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	name = "given post when requested then return 200"
	t.Run(name, func(t *testing.T) {
		// Arrange
		post, err := postFactory.CreatePost()
		assert.NoError(t, err)
		url := fmt.Sprintf("/v1/posts/%d", post.ID)

		// Act
		rec, _ := server.Get(t, url, nil)

		// Assert
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}