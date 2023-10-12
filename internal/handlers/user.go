package handlers

import (
	"errors"
	"net/http"

	"nexablog/internal/models"
	"nexablog/internal/services"
	"nexablog/internal/services/permission"
	"nexablog/internal/services/post"
	"nexablog/internal/services/user"
	"nexablog/internal/utils"
	"nexablog/pkg/lib"
	"nexablog/pkg/validator"
)

type User struct {
	UserSvc       user.Service
	PermissionSvc permission.Service
	PostSvc       post.Service
}

func (h *User) Register(w http.ResponseWriter, r *http.Request) error {
	payload := models.UserIn{}

	if err := utils.ReadJson(w, r, &payload); err != nil {
		return utils.NewApiError(err.Error(), http.StatusUnprocessableEntity)
	}

	v := validator.New()

	if models.ValidateUser(v, &payload); !v.Valid() {
		return utils.WriteJson(w, http.StatusUnprocessableEntity, lib.H[any]{
			"details": v.Errors,
		})
	}

	user, err := h.UserSvc.CreateUser(r.Context(), payload)

	if errors.Is(err, services.ErrDuplicateKey) {
		return utils.NewApiError("email already exists", http.StatusBadGateway)
	}

	if err != nil {
		return err
	}

	err = h.PermissionSvc.AddUserPermission(
		r.Context(),
		user.UserID,
		"posts:read",
		"posts:write",
	)

	if err != nil {
		return err
	}

	return utils.WriteJson(w, http.StatusCreated, user)
}

func (h *User) GetMe(w http.ResponseWriter, r *http.Request) error {
	user := utils.GetUser(r)

	posts, err := h.PostSvc.FindPostsByAuthor(r.Context(), user.UserID)
	if err != nil {
		return err
	}

	return utils.WriteJson(w, http.StatusOK, lib.H[any]{
		"user":  user,
		"posts": posts,
	})
}
