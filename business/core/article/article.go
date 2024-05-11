package article

import (
	"context"
	"errors"
	"github/islamghany/blog/business/core/user"
	"github/islamghany/blog/foundation/logger"
	"time"
)

var (
	ErrorNotFound     = errors.New("article not found")
	ErrUserIsDisabled = errors.New("user is disabled")
)

type Storer interface {
	Create(ctx context.Context, art Article) (int, error)
	QueryByID(ctx context.Context, id int) (Article, error)
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
