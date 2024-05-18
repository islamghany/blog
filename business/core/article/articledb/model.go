package articledb

import (
	"github/islamghany/blog/business/core/article"
	"github/islamghany/blog/business/core/user/userdb"
	"github/islamghany/blog/business/data/dbsql/pgx/dbarray"
	"time"

	"github.com/google/uuid"
)

type dbarticle struct {
	ID        int            `db:"id"`
	Title     string         `db:"title"`
	Content   string         `db:"content"`
	Tags      dbarray.String `db:"tags"`
	AuthorID  uuid.UUID      `db:"author_id"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func toDBArticle(art article.Article) dbarticle {
	return dbarticle{
		ID:        art.ID,
		Title:     art.Title,
		Content:   art.Content,
		Tags:      art.Tags,
		AuthorID:  art.AuthorID,
		CreatedAt: art.CreatedAt.UTC(),
		UpdatedAt: art.UpdatedAt.UTC(),
	}
}

func toArticle(dbart dbarticle) article.Article {
	return article.Article{
		ID:        dbart.ID,
		Title:     dbart.Title,
		Content:   dbart.Content,
		Tags:      dbart.Tags,
		AuthorID:  dbart.AuthorID,
		CreatedAt: dbart.CreatedAt,
		UpdatedAt: dbart.UpdatedAt,
	}
}

func toArticleSlice(dbarts []dbarticle) []article.Article {
	arts := make([]article.Article, len(dbarts))
	for i, dbart := range dbarts {
		arts[i] = toArticle(dbart)
	}
	return arts
}

func toDBArticleSlice(arts []article.Article) []dbarticle {
	dbarts := make([]dbarticle, len(arts))
	for i, art := range arts {
		dbarts[i] = toDBArticle(art)
	}
	return dbarts
}

// ArticleWithAuthor
type DBArticleWithAuthor struct {
	ID        int            `db:"pid"`
	Title     string         `db:"title"`
	Content   string         `db:"content"`
	Tags      dbarray.String `db:"tags"`
	CreatedAt time.Time      `db:"pcreated_at"`
	UpdatedAt time.Time      `db:"pupdated_at"`
	userdb.DBUser
}

func toDBArticleWithAuthor(art article.ArticleWithAuthor) DBArticleWithAuthor {
	return DBArticleWithAuthor{
		ID:        art.ID,
		Title:     art.Title,
		Content:   art.Content,
		Tags:      art.Tags,
		CreatedAt: art.CreatedAt,
		UpdatedAt: art.UpdatedAt,
		DBUser:    userdb.ToDBUser(art.User),
	}
}

func toArticleWithAuthor(dbart DBArticleWithAuthor) article.ArticleWithAuthor {
	return article.ArticleWithAuthor{
		ID:        dbart.ID,
		Title:     dbart.Title,
		Content:   dbart.Content,
		Tags:      dbart.Tags,
		CreatedAt: dbart.CreatedAt,
		UpdatedAt: dbart.UpdatedAt,
		User:      userdb.ToCoreUser(dbart.DBUser),
	}
}

func toArticleWithAuthorSlice(dbarts []DBArticleWithAuthor) []article.ArticleWithAuthor {
	arts := make([]article.ArticleWithAuthor, len(dbarts))
	for i, dbart := range dbarts {
		arts[i] = toArticleWithAuthor(dbart)
	}
	return arts
}
