package handlers

import (
	"net/http"

	"nexablog/internal/utils"
	"nexablog/pkg/lib"
)

func Welcome(w http.ResponseWriter, r *http.Request) error {
	return utils.WriteJson(w, http.StatusOK, lib.H[string]{
		"detail": "welcome to nexablog",
	})
}
