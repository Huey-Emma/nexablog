package utils

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"nexablog/internal/models"
)

type DBTX interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

type Row interface {
	Scan(...any) error
}

type ApiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	errmsg string
	code   int
}

func (e ApiError) Error() string {
	return e.errmsg
}

func (e ApiError) Code() int {
	return e.code
}

func NewApiError(errmsg string, code int) *ApiError {
	return &ApiError{errmsg, code}
}

func SendStatus(w http.ResponseWriter, code int) error {
	w.WriteHeader(code)
	return nil
}

func WriteJson(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func ReadJson(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, int64(1_048_576))

	defer func() {
		_ = r.Body.Close()
	}()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(v)

	var syntaxErr *json.SyntaxError

	if errors.As(err, &syntaxErr) {
		return fmt.Errorf("malformed json at %d", syntaxErr.Offset)
	}

	if errors.Is(err, io.EOF) {
		return fmt.Errorf("request body has no content")
	}

	if errors.Is(err, io.EOF) {
		return fmt.Errorf("malformed json")
	}

	if err != nil {
		return err
	}

	return nil
}

func GetPasswordHash(plain string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, err
	}

	return hash, nil
}

func PasswordMatch(hash []byte, plain string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(plain))

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func GenRandString(n int) (string, error) {
	rb := make([]byte, n)

	_, err := rand.Read(rb)
	if err != nil {
		return "", err
	}

	encoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	return encoder.EncodeToString(rb), nil
}

func HashRandString(s string) []byte {
	h := sha256.Sum256([]byte(s))
	return h[:]
}

func SetUser(r *http.Request, u *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), "user", u)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value("user").(*models.User)

	if !ok || user == nil {
		panic("no user in request context")
	}

	return user
}
