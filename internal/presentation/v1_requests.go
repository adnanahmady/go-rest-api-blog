package presentation

import (
	"github.com/adnanahmady/go-rest-api-blog/pkg/errs"
)

type SavePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (r *SavePostRequest) Validate() error {
	if r.Title == "" {
		return errs.NewValidationError(errs.Errors{
			"title": "title is required",
		})
	}
	if r.Content == "" {
		return errs.NewValidationError(errs.Errors{
			"content": "content is required",
		})
	}
	return nil
}
