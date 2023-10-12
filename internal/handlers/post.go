package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"nexablog/internal/models"
	"nexablog/internal/services"
	"nexablog/internal/services/post"
	"nexablog/internal/utils"
	"nexablog/pkg/lib"
	"nexablog/pkg/validator"
)

type Post struct {
	PostSvc post.Service
}

func (h *Post) CreatePost(w http.ResponseWriter, r *http.Request) error {
	payload := models.PostIn{}

	if err := utils.ReadJson(w, r, &payload); err != nil {
		return utils.NewApiError(err.Error(), http.StatusUnprocessableEntity)
	}

	v := validator.New()

	if models.ValidatePost(v, &payload); !v.Valid() {
		return utils.WriteJson(w, http.StatusUnprocessableEntity, lib.H[any]{
			"details": v.Errors,
		})
	}

	payload.AuthorID = utils.GetUser(r).UserID

	post, err := h.PostSvc.CreatePost(r.Context(), payload)
	if err != nil {
		return err
	}

	return utils.WriteJson(w, http.StatusCreated, post)
}

func (h *Post) FindAllPosts(w http.ResponseWriter, r *http.Request) error {
	posts, err := h.PostSvc.FindAllPosts(r.Context())
	if err != nil {
		return err
	}

	return utils.WriteJson(w, http.StatusOK, posts)
}

func (h *Post) FindPostByID(w http.ResponseWriter, r *http.Request) error {
	postID, _ := strconv.Atoi(chi.URLParam(r, "post-id"))

	post, err := h.PostSvc.FindPostByID(r.Context(), postID)

	if errors.Is(err, services.ErrResourceNotFound) {
		return utils.WriteJson(w, http.StatusNotFound, lib.H[string]{
			"detail": "post not found",
		})
	}

	if err != nil {
		return err
	}

	return utils.WriteJson(w, http.StatusOK, post)
}

func (h *Post) DeletePostByID(w http.ResponseWriter, r *http.Request) error {
	postID, _ := strconv.Atoi(chi.URLParam(r, "post-id"))

	post, err := h.PostSvc.FindPostByID(r.Context(), postID)

	if errors.Is(err, services.ErrResourceNotFound) {
		return utils.NewApiError("post not found", http.StatusNotFound)
	}

	if err != nil {
		return err
	}

	user := utils.GetUser(r)

	if !user.IsOwner(post.AuthorID) {
		return utils.NewApiError("not allowed", http.StatusForbidden)
	}

	err = h.PostSvc.DeletePostByID(r.Context(), postID)

	if err != nil && !errors.Is(err, services.ErrResourceNotFound) {
		return err
	}

	return utils.SendStatus(w, http.StatusNoContent)
}

func (h *Post) UpdatePostByID(w http.ResponseWriter, r *http.Request) error {
	postID, _ := strconv.Atoi(chi.URLParam(r, "post-id"))

	post, err := h.PostSvc.FindPostByID(r.Context(), postID)

	if errors.Is(err, services.ErrResourceNotFound) {
		return utils.NewApiError("post not found", http.StatusNotFound)
	}

	if err != nil {
		return err
	}

	user := utils.GetUser(r)

	if !user.IsOwner(post.AuthorID) {
		return utils.NewApiError("not allowed", http.StatusForbidden)
	}

	payload := models.PostIn{}

	if err := utils.ReadJson(w, r, &payload); err != nil {
		return utils.NewApiError(err.Error(), http.StatusUnprocessableEntity)
	}

	v := validator.New()

	if models.ValidatePost(v, &payload); !v.Valid() {
		return utils.WriteJson(w, http.StatusUnprocessableEntity, lib.H[any]{
			"details": v.Errors,
		})
	}

	post, err = h.PostSvc.UpdatePostByID(
		r.Context(),
		payload,
		post.PostID,
		post.Version,
	)

	if err != nil && !errors.Is(err, services.ErrResourceNotFound) {
		return err
	}

	return utils.WriteJson(w, http.StatusOK, post)
}
