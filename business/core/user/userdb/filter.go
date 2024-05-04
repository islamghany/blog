package userdb

import (
	"bytes"
	"github/islamghany/blog/business/core/user"
	"strings"
)

func (s *Store) applyFilters(filters user.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	// where clause
	var wc []string

	if filters.ID != nil {
		data["id"] = *filters.ID
		wc = append(wc, "id = :id")
	}
	if filters.Username != nil {
		data["username"] = *filters.Username
		wc = append(wc, "username = :username")
	}
	if filters.Email != nil {
		data["email"] = *filters.Email
		wc = append(wc, "email = :email")
	}
	if filters.FirstName != nil {
		data["first_name"] = *filters.FirstName
		wc = append(wc, "first_name = :first_name")
	}
	if filters.LastName != nil {
		data["last_name"] = *filters.LastName
		wc = append(wc, "last_name = :last_name")
	}
	if filters.CreatedAt != nil {
		data["created_at"] = *filters.CreatedAt
		wc = append(wc, "created_at = :created_at")
	}
	if filters.UpdatedAt != nil {
		data["updated_at"] = *filters.UpdatedAt
		wc = append(wc, "updated_at = :updated_at")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
