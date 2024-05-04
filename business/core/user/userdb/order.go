package userdb

import (
	"fmt"
	"github/islamghany/blog/business/core/user"
	"github/islamghany/blog/business/data/order"
)

var OrderByFields = map[string]string{
	user.OrderByID:        "id",
	user.OrderByEmail:     "email",
	user.OrderByUsername:  "username",
	user.OrderByFirstName: "first_name",
	user.OrderByLastName:  "last_name",
	user.OrderByCreatedAt: "created_at",
	user.OrderByUpdatedAt: "updated_at",
}

// ApplySort applies the sorting options for users.
func orderByClause(orderBy order.By) (string, error) {
	by, exists := OrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("invalid field %q", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Dir, nil
}
