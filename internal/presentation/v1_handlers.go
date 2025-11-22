package presentation

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/adnanahmady/go-rest-api-blog/internal/application"
	"github.com/adnanahmady/go-rest-api-blog/pkg/errs"
	"github.com/adnanahmady/go-rest-api-blog/pkg/request"
	"github.com/adnanahmady/go-rest-api-blog/pkg/response"
	"github.com/go-chi/chi/v5"
)

type V1Handlers struct {
	create application.CreatePostUseCase
	list   application.ListPostsUseCase
	show   application.ShowPostUseCase
	update application.UpdatePostUseCase
	delete application.DeletePostUseCase
}

func NewV1Handlers(
	create application.CreatePostUseCase,
	list application.ListPostsUseCase,
	show application.ShowPostUseCase,
	update application.UpdatePostUseCase,
	delete application.DeletePostUseCase,
) *V1Handlers {
	return &V1Handlers{
		create: create,
		list:   list,
		show:   show,
		update: update,
		delete: delete,
	}
}

// @Summary Health check
// @Description Health check
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} response.JsonResponse{data=string}
// @Failure 500 {object} errs.AppError
// @Router /health [get]
func (h *V1Handlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// @Summary Create a new post
// @Description Create a new post
// @Tags posts
// @Accept json
// @Produce json
// @Param post body SavePostRequest true "Post data"
// @Success 201 {object} response.JsonResponse{data=PostResource}
// @Failure 400 {object} errs.AppError
// @Failure 422 {object} errs.AppError
// @Failure 500 {object} errs.AppError
// @Router /posts [post]
func (h *V1Handlers) CreatePost(w http.ResponseWriter, r *http.Request) {
	lgr := request.GetLogger(r.Context())

	var req SavePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lgr.Error("failed to decode request body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		lgr.Error("request validation failed", err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	post, err := h.create.Create(r.Context(), req.Title, req.Content)
	if err != nil {
		var appErr *errs.AppError
		if errors.As(err, &appErr) {
			lgr.Error("failed to create post", appErr)
			http.Error(w, appErr.Error(), appErr.Code)
			return
		}

		lgr.Error("failed to create post", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	res := newPostResource(*post)
	json.NewEncoder(w).Encode(response.NewJsonResponse(res, nil))
}

// @Summary List all posts
// @Description List all posts
// @Tags posts
// @Accept json
// @Produce json
// @Param page query int true "Page number"
// @Param per_page query int true "Per page"
// @Success 200 {object} response.JsonResponse{data=PostResource,meta=map[string]any{pagination=application.PaginationDTO}}
// @Failure 400 {object} errs.AppError
// @Failure 500 {object} errs.AppError
// @Router /posts [get]
func (h *V1Handlers) ListPosts(w http.ResponseWriter, r *http.Request) {
	lgr := request.GetLogger(r.Context())

	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 32)
	if err != nil {
		page = 1
	}

	perPage, err := strconv.ParseInt(r.URL.Query().Get("per_page"), 10, 32)
	if err != nil {
		perPage = 10
	}

	dto, err := h.list.List(r.Context(), int(page), int(perPage))
	if err != nil {
		var appErr *errs.AppError
		if errors.As(err, &appErr) {
			lgr.Error("failed to get posts list", appErr)
			http.Error(w, appErr.Error(), appErr.Code)
			return
		}
		lgr.Error("failed to get posts list", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	collection := newPostCollection(dto.Data)
	res := response.NewJsonResponse(collection, map[string]any{
		"pagination": dto,
	})
	json.NewEncoder(w).Encode(res)
}

// @Summary Show a post
// @Description Show a post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} response.JsonResponse{data=PostResource}
// @Failure 400 {object} errs.AppError
// @Failure 404 {object} errs.AppError
// @Failure 500 {object} errs.AppError
// @Router /posts/{id} [get]
func (h *V1Handlers) ShowPost(w http.ResponseWriter, r *http.Request) {
	lgr := request.GetLogger(r.Context())

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		lgr.Error("failed to parse post id", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := h.show.Show(r.Context(), uint(id))
	if err != nil {
		var appErr *errs.AppError
		if errors.As(err, &appErr) {
			lgr.Error("post not found", appErr)
			http.Error(w, appErr.Error(), appErr.Code)
			return
		}
		lgr.Error("failed to get post", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	res := newPostResource(*post)
	json.NewEncoder(w).Encode(response.NewJsonResponse(res, nil))
}

// @Summary Update a post
// @Description Update a post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param post body SavePostRequest true "Post data"
// @Success 200 {object} response.JsonResponse{data=PostResource}
// @Failure 400 {object} errs.AppError
// @Failure 422 {object} errs.AppError
// @Failure 500 {object} errs.AppError
// @Router /posts/{id} [put]
func (h *V1Handlers) UpdatePost(w http.ResponseWriter, r *http.Request) {
	lgr := request.GetLogger(r.Context())

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		lgr.Error("failed to parse post id", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req SavePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lgr.Error("failed to decode request body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		lgr.Error("request validation failed", err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	post, err := h.update.Update(r.Context(), uint(id), req.Title, req.Content)
	if err != nil {
		var appErr *errs.AppError
		if errors.As(err, &appErr) {
			lgr.Error("post not found", appErr)
			http.Error(w, appErr.Error(), appErr.Code)
			return
		}
		lgr.Error("failed to update post", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	res := newPostResource(*post)
	json.NewEncoder(w).Encode(response.NewJsonResponse(res, nil))
}

// @Summary Delete a post
// @Description Delete a post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 204 {object} response.JsonResponse{data=nil}
// @Failure 400 {object} errs.AppError
// @Failure 404 {object} errs.AppError
// @Failure 500 {object} errs.AppError
// @Router /posts/{id} [delete]
func (h *V1Handlers) DeletePost(w http.ResponseWriter, r *http.Request) {
	lgr := request.GetLogger(r.Context())

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		lgr.Error("failed to parse post id", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.delete.Delete(r.Context(), uint(id)); err != nil {
		var appErr *errs.AppError
		if errors.As(err, &appErr) {
			lgr.Error("post not found", appErr)
			http.Error(w, appErr.Error(), appErr.Code)
			return
		}
		lgr.Error("failed to delete post", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
