package article

import (
	"time"

	"github.com/google/uuid"
)

// Article represents a blog post.
type Article struct {
	ID        int
	Title     string
	Content   string
	Tags      []string
	AuthorID  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewArticle represents data used to create an article.
type NewArticle struct {
	Title    string
	Content  string
	Tags     []string
	AuthorID uuid.UUID
}

// UpdateArticle represents data used to update an article.
type UpdateArticle struct {
	Title   *string
	Content *string
	Tags    []string
}
