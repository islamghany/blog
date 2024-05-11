package articleapi

import (
	"github/islamghany/blog/business/core/article"
	"github/islamghany/blog/foundation/validate"
	"time"

	"github.com/google/uuid"
)

type ApiArticle struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	AuthorID  uuid.UUID `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toApiArticle(a article.Article) ApiArticle {
	return ApiArticle{
		ID:        a.ID,
		Title:     a.Title,
		Content:   a.Content,
		Tags:      a.Tags,
		AuthorID:  a.AuthorID,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

type ApiNewArticle struct {
	Title    string    `json:"title" validate:"required,min=3,max=100"`
	Content  string    `json:"content" validate:"required,min=3,max=1500"`
	Tags     []string  `json:"tags" validate:"required,min=1"`
	AuthorID uuid.UUID `json:"author_id" validate:"required,uuid"`
}

func (a *ApiNewArticle) Validate() error {
	return validate.Check(a)
}

func toNewArticleCore(a ApiNewArticle) article.NewArticle {
	return article.NewArticle{
		Title:    a.Title,
		Content:  a.Content,
		Tags:     a.Tags,
		AuthorID: a.AuthorID,
	}
}
