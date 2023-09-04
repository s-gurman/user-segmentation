package domain

import "errors"

var ErrEmptySlug = errors.New("empty segment name")

type Segment struct {
	ID   int  `db:"id"`
	Slug Slug `db:"slug"`
}

type Slug string

func NewSlug(name string) (Slug, error) {
	if name == "" {
		return Slug(""), ErrEmptySlug
	}
	return Slug(name), nil
}

func (slug Slug) String() string {
	return string(slug)
}
