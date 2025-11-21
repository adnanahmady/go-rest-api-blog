package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/adnanahmady/go-rest-api-blog/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestUpdatePost(t *testing.T) {
	server, err := test.Setup()
	assert.NoError(t, err)
	defer server.Close()
	postFactory := NewPostFactory(server.App.Database)
	type testData struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	name := "given data when its valid then it should update the post"
	t.Run(name, func(t *testing.T) {
		// Arrange
		post, err := postFactory.CreatePost()
		assert.NoError(t, err)
		url := fmt.Sprintf("/v1/posts/%d", post.ID)
		data := testData{
			Title:   "Test Post",
			Content: "Test Content",
		}

		// Act
		_, body := server.Put(t, url, data, nil)

		// Assert
		bodyMap := body.(map[string]any)
		dataMap := bodyMap["data"].(map[string]any)

		assert.IsType(t, float64(0), dataMap["id"])
		assert.Equal(t, data.Title, dataMap["title"])
		assert.Equal(t, data.Content, dataMap["content"])
		assert.IsType(t, "string", dataMap["created_at"])
		assert.IsType(t, "string", dataMap["updated_at"])
		test.AssertTimeFormat(t, dataMap["created_at"].(string))
		test.AssertTimeFormat(t, dataMap["updated_at"].(string))
	})

	name = "given data when its valid then return 201"
	t.Run(name, func(t *testing.T) {
		// Arrange
		post, err := postFactory.CreatePost()
		assert.NoError(t, err)
		url := fmt.Sprintf("/v1/posts/%d", post.ID)
		data := testData{
			Title:   "Test Post",
			Content: "Test Content",
		}

		// Act
		rec, _ := server.Put(t, url, data, nil)

		// Assert
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	vals := []struct {
		name string
		data testData
		exp  int
		errs map[string]any
	}{
		{
			name: "given data when title is empty then return error",
			data: testData{
				Title:   "",
				Content: "Test Content",
			},
			exp: http.StatusUnprocessableEntity,
			errs: map[string]any{
				"title": "title is required",
			},
		},
		{
			name: "given data when content is empty then return error",
			data: testData{
				Title:   "Test Post",
				Content: "",
			},
			exp: http.StatusUnprocessableEntity,
			errs: map[string]any{
				"content": "content is required",
			},
		},
	}
	for _, val := range vals {
		t.Run(val.name, func(t *testing.T) {
			// Arrange
			post, err := postFactory.CreatePost()
			assert.NoError(t, err)
			url := fmt.Sprintf("/v1/posts/%d", post.ID)

			// Act
			rec, body := server.Put(t, url, val.data, nil)

			// Assert
			assert.Equal(t, val.exp, rec.Code)
			assert.Equal(t, val.errs, body.(map[string]any)["errors"])
		})
	}
}
