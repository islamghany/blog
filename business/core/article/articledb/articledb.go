package articledb

import (
	"context"
	"errors"
	"fmt"
	"github/islamghany/blog/business/core/article"
	db "github/islamghany/blog/business/data/dbsql/pgx"
	"github/islamghany/blog/foundation/logger"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	Log *logger.Logger
	DB  *sqlx.DB
}

func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		Log: log,
		DB:  db,
	}
}

func (s *Store) Create(ctx context.Context, art article.Article) (int, error) {
	q := `
		INSERT INTO articles 
			(title, content, tags, author_id, created_at, updated_at)
		VALUES
			(:title, :content, :tags, :author_id, :created_at, :updated_at)
		RETURNING id
	 `
	dbart := toDBArticle(art)
	ret := struct {
		ID int `db:"id"`
	}{}
	if err := db.NamedQueryStruct(ctx, s.Log, s.DB, q, dbart, &ret); err != nil {
		return 0, err
	}
	return ret.ID, nil
}

func (s *Store) QueryByID(ctx context.Context, id int) (article.Article, error) {
	args := struct {
		ID int `db:"id"`
	}{ID: id}
	q := `
		SELECT * FROM articles WHERE id = :id;
	`
	dbart := dbarticle{}
	if err := db.NamedQueryStruct(ctx, s.Log, s.DB, q, args, &dbart); err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return article.Article{}, article.ErrorNotFound
		}
		return article.Article{}, fmt.Errorf("querying for article %d: %w", id, err)
	}
	if dbart.ID == 0 {
		return article.Article{}, article.ErrorNotFound
	}
	return toArticle(dbart), nil
}

func (s *Store) Update(ctx context.Context, art article.Article) error {
	q := `
	 UPDATE articles
	 SET
		title = :title,
		content = :content,
		tags = :tags,
		updated_at = :updated_at
	 WHERE id = :id and author_id = :author_id
   `

	dbart := toDBArticle(art)
	if err := db.NamedExecContext(ctx, s.Log, s.DB, q, dbart); err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return article.ErrorNotFound
		}
		return fmt.Errorf("updating article %d: %w", art.ID, err)
	}
	return nil
}
