package article

import (
	"context"
	"errors"
	"github/islamghany/blog/business/core/user"
	"github/islamghany/blog/business/data/transaction"
	"github/islamghany/blog/foundation/logger"
	"time"
)

var (
	ErrorNotFound     = errors.New("article not found")
	ErrUserIsDisabled = errors.New("user is disabled")
)

type Storer interface {
	ExecuteUnderTransaction(tx transaction.Transaction) (Storer, error)
	Create(ctx context.Context, art Article) (int, error)
	QueryByID(ctx context.Context, id int) (Article, error)
	Update(ctx context.Context, art Article) error
	Query(ctx context.Context, search string, pageNumber, rowsPerPage int) ([]ArticleWithAuthor, int, error)
}

type Core struct {
	store    Storer
	log      *logger.Logger
	userCore *user.Core
}

func NewCore(log *logger.Logger, store Storer, userCore *user.Core) *Core {
	return &Core{
		store:    store,
		log:      log,
		userCore: userCore,
	}
}

// ExecuteUnderTransaction constructs a new Core value that will use the
// specified transaction in any store related calls.
func (c *Core) ExecuteUnderTransaction(tx transaction.Transaction) (*Core, error) {
	trS, err := c.store.ExecuteUnderTransaction(tx)
	if err != nil {
		return nil, err
	}

	c = &Core{
		store:    trS,
		log:      c.log,
		userCore: c.userCore,
	}

	return c, nil
}

func (c *Core) Create(ctx context.Context, na NewArticle) (int, error) {
	now := time.Now()
	art := Article{
		Title:     na.Title,
		Content:   na.Content,
		Tags:      na.Tags,
		AuthorID:  na.AuthorID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	id, err := c.store.Create(ctx, art)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *Core) QueryByID(ctx context.Context, id int) (Article, error) {
	art, err := c.store.QueryByID(ctx, id)
	if err != nil {
		return Article{}, err
	}
	return art, nil
}

func (c *Core) Update(ctx context.Context, art Article, ua UpdateArticle) error {
	if ua.Title != nil {
		art.Title = *ua.Title
	}
	if ua.Content != nil {
		art.Content = *ua.Content
	}
	if len(ua.Tags) > 0 {
		art.Tags = ua.Tags
	}
	art.UpdatedAt = time.Now().UTC()
	err := c.store.Update(ctx, art)
	if err != nil {
		return err
	}
	return nil
}

func (c *Core) Query(ctx context.Context, search string, pageNumber, rowsPerPage int) ([]ArticleWithAuthor, int, error) {
	return c.store.Query(ctx, search, pageNumber, rowsPerPage)
}
