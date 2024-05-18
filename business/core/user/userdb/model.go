package userdb

import (
	"github/islamghany/blog/business/core/user"
	"github/islamghany/blog/business/data/dbsql/pgx/dbarray"
	"time"

	"github.com/google/uuid"
)

type DBUser struct {
	ID             uuid.UUID      `db:"id"`
	Email          string         `db:"email"`
	Username       string         `db:"username"`
	Roles          dbarray.String `db:"roles"`
	FirstName      string         `db:"first_name"`
	LastName       string         `db:"last_name"`
	PasswordHashed []byte         `db:"password_hashed"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
	Enabled        bool           `db:"enabled"`
	Version        int            `db:"version"`
}

func ToDBUser(user user.User) DBUser {
	roles := make(dbarray.String, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name()
	}
	return DBUser{
		ID:             user.ID,
		Email:          user.Email,
		Username:       user.Username,
		Roles:          roles,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		PasswordHashed: user.PasswordHashed,
		CreatedAt:      user.CreatedAt.UTC(),
		UpdatedAt:      user.UpdatedAt.UTC(),
		Enabled:        user.Enabled,
		Version:        user.Version,
	}
}

func ToCoreUser(dbusr DBUser) user.User {
	roles := make([]user.Role, len(dbusr.Roles))
	for i, role := range dbusr.Roles {
		r, err := user.ParseRole(role)
		if err != nil {
			panic(err)
		}
		roles[i] = r
	}
	return user.User{
		ID:             dbusr.ID,
		Email:          dbusr.Email,
		Username:       dbusr.Username,
		Roles:          roles,
		FirstName:      dbusr.FirstName,
		LastName:       dbusr.LastName,
		PasswordHashed: dbusr.PasswordHashed,
		CreatedAt:      dbusr.CreatedAt,
		UpdatedAt:      dbusr.UpdatedAt,
		Enabled:        dbusr.Enabled,
		Version:        dbusr.Version,
	}

}

func ToCoreUserSlice(DBUsers []DBUser) []user.User {
	users := make([]user.User, len(DBUsers))
	for i, DBUser := range DBUsers {
		usr := ToCoreUser(DBUser)
		users[i] = usr
	}
	return users
}

type DBUserWithTotal struct {
	Total int `db:"total"`
	DBUser
}

func ToCoreUserWithTotal(DBUser DBUserWithTotal) (user.User, int) {
	usr := ToCoreUser(DBUser.DBUser)
	return usr, DBUser.Total
}

func ToCoreUserWithTotalSlice(DBUsers []DBUserWithTotal) ([]user.User, int) {
	users := make([]user.User, len(DBUsers))
	var total int
	for i, DBUser := range DBUsers {
		usr, t := ToCoreUserWithTotal(DBUser)
		users[i] = usr
		total = t
	}
	return users, total
}
