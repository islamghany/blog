package mid

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github/islamghany/blog/business/data/transaction"
	"github/islamghany/blog/foundation/logger"
	"github/islamghany/blog/foundation/web"
	"net/http"
)

func ExecuteInTransaction(log *logger.Logger, bgn transaction.Beginner) web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hasCommited := false
			log.Info(ctx, "Begin transaction")
			tx, err := bgn.Begin()
			if err != nil {
				return fmt.Errorf("BEGIN TRANSACTION: %w", err)
			}
			defer func() {
				if !hasCommited {
					log.Info(ctx, "Rollback transaction")
				}
				if err := tx.Rollback(); err != nil {
					if errors.Is(err, sql.ErrTxDone) {
						return
					}
					log.Error(ctx, "Rollback transaction", err)
				}
			}()
			ctx = transaction.Set(ctx, tx)
			if err := handler(ctx, w, r); err != nil {
				return fmt.Errorf("EXECUTE TRANSACTION: %w", err)
			}
			log.Info(ctx, "COMMIT TRANSACTION")
			if err := tx.Commit(); err != nil {
				return fmt.Errorf("COMMIT TRANSACTION: %w", err)
			}
			hasCommited = true
			return nil
		}

	}
}
