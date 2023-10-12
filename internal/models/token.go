package models

import "time"

type Scope string

const (
	ScopeAuthentication Scope = "authentication"
)

type Token struct {
	Hash      []byte
	UserID    int
	ExpiresAt time.Time
	Scope     Scope
}

type TokenIn struct {
	UserID    int
	ExpiresAt time.Time
	Scope     Scope
}

type TokenOut struct {
	Plain     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
