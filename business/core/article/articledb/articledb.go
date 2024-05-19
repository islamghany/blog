package articledb

import (
	"context"
	"errors"
	"fmt"
	"github/islamghany/blog/business/core/article"
	db "github/islamghany/blog/business/data/dbsql/pgx"
	"github/islamghany/blog/business/data/transaction"
	"github/islamghany/blog/foundation/logger"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	Log *logger.Logger
	DB  sqlx.ExtContext
}

func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		Log: log,
		DB:  db,
	}
}
func (s *Store) ExecuteUnderTransaction(tx transaction.Transaction) (article.Storer, error) {
	ec, err := db.GetExtContext(tx)
	if err != nil {
		return nil, err
	}

	as := &Store{
		Log: s.Log,
		DB:  ec,
	}

	return as, nil
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

func (s *Store) Query(ctx context.Context, search string, pageNumber, rowsPerPage int) ([]article.ArticleWithAuthor, int, error) {
	data := map[string]any{
		"offset": (pageNumber - 1) * rowsPerPage,
		"limit":  rowsPerPage,
		"search": search,
	}
	q := `
		SELECT 
			 a.id as pid, title, substring(content, 0, 100) as content, tags, a.created_at as pcreated_at, a.updated_at as pupdated_at, 
			 username, email, u.created_at as created_at, u.updated_at as updated_at, u.id as id
		FROM articles a
		JOIN users u ON a.author_id = u.id
		WHERE (a.tsv_document @@ to_tsquery('english', :search))
		ORDER BY ts_rank(a.tsv_document, to_tsquery('english', :search)) DESC
		OFFSET :offset
		LIMIT :limit
	`
	var dbarts []DBArticleWithAuthor
	if err := db.NamedQuerySlice(ctx, s.Log, s.DB, q, data, &dbarts); err != nil {
		return nil, 0, fmt.Errorf("querying articles: %w", err)
	}
	return toArticleWithAuthorSlice(dbarts), len(dbarts), nil
}
