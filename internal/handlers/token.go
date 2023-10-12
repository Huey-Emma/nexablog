package handlers

import (
	"errors"
	"net/http"
	"time"

	"nexablog/internal/models"
	"nexablog/internal/services"
	"nexablog/internal/services/token"
	"nexablog/internal/services/user"
	"nexablog/internal/utils"
)

type Token struct {
	UserSvc  user.Service
	TokenSvc token.Service
}

func (h *Token) GetToken(w http.ResponseWriter, r *http.Request) error {
	payload := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})

	if err := utils.ReadJson(w, r, payload); err != nil {
		return utils.NewApiError(err.Error(), http.StatusUnprocessableEntity)
	}

	user, err := h.UserSvc.FindUserByEmail(r.Context(), payload.Email)

	if errors.Is(err, services.ErrResourceNotFound) {
		return utils.NewApiError("invalid credentials", http.StatusUnauthorized)
	}

	if err != nil {
		return err
	}

	isMatch, err := utils.PasswordMatch(user.Password, payload.Password)
	if err != nil {
		return err
	}

	if !isMatch {
		return utils.NewApiError("invalid credentials", http.StatusUnauthorized)
	}

	token, err := h.TokenSvc.AddToken(r.Context(), models.TokenIn{
		UserID:    user.UserID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Scope:     models.ScopeAuthentication,
	})
	if err != nil {
		return err
	}

	return utils.WriteJson(w, http.StatusCreated, token)
}
