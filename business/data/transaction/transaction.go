package transaction

import (
	"context"
)

type Transaction interface {
	Commit() error
	Rollback() error
}

type Beginner interface {
	Begin() (Transaction, error)
}

type transactionKey int

const key transactionKey = 1

func Set(ctx context.Context, tx Transaction) context.Context {
	return context.WithValue(ctx, key, tx)
}

func Get(ctx context.Context) (Transaction, bool) {
	v, ok := ctx.Value(key).(Transaction)
	return v, ok
}
