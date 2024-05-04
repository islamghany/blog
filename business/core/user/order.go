package user

import "github/islamghany/blog/business/data/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByID, order.ASC)

// Set of fields that the results can be ordered by. These are the names
// that should be used by the application layer.
const (
	OrderByID        = "id"
	OrderByEmail     = "email"
	OrderByUsername  = "username"
	OrderByFirstName = "first_name"
	OrderByLastName  = "last_name"
	OrderByCreatedAt = "created_at"
	OrderByUpdatedAt = "updated_at"
)
