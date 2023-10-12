package models

import (
	"time"

	"nexablog/pkg/lib"
	"nexablog/pkg/validator"
)

type Post struct {
	PostID    int       `json:"post_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	AuthorID  int       `json:"-"`
	Version   int       `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type Posts []Post

type PostIn struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	AuthorID int    `json:"-"`
}

func ValidatePost(v *validator.Validator, p *PostIn) {
	v.Check(lib.NonWhiteSpace(p.Title), "title", "cannot be blank")
	v.Check(lib.NonWhiteSpace(p.Body), "body", "cannot be blank")
}
