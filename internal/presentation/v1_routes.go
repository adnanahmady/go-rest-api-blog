package presentation

import (
	"github.com/adnanahmady/go-rest-api-blog/pkg/request"
	"github.com/go-chi/chi/v5"
)

type V1Routes struct{}

func NewV1Routes(router request.Router, v1Handlers *V1Handlers) *V1Routes {
	router.GetEngine().Route("/v1", func(r chi.Router) {
		r.Get("/health", v1Handlers.HealthCheck)

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", v1Handlers.CreatePost)
			r.Get("/", v1Handlers.ListPosts)
			r.Get("/{id}", v1Handlers.ShowPost)
			r.Put("/{id}", v1Handlers.UpdatePost)
			r.Delete("/{id}", v1Handlers.DeletePost)
		})
	})

	return &V1Routes{}
}
