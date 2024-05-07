package userdb

import (
	"github/islamghany/blog/business/core/user"
	"github/islamghany/blog/business/data/dbsql/pgx/dbarray"
	"time"

	"github.com/google/uuid"
)

type dbuser struct {
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

func toDBUser(user user.User) dbuser {
	roles := make(dbarray.String, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name()
	}
	return dbuser{
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

func toCoreUser(dbusr dbuser) user.User {
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

func toCoreUserSlice(dbusers []dbuser) []user.User {
	users := make([]user.User, len(dbusers))
	for i, dbuser := range dbusers {
		usr := toCoreUser(dbuser)
		users[i] = usr
	}
	return users
}

type dbuserWithTotal struct {
	Total int `db:"total"`
	dbuser
}

func toCoreUserWithTotal(dbuser dbuserWithTotal) (user.User, int) {
	usr := toCoreUser(dbuser.dbuser)
	return usr, dbuser.Total
}

func toCoreUserWithTotalSlice(dbusers []dbuserWithTotal) ([]user.User, int) {
	users := make([]user.User, len(dbusers))
	var total int
	for i, dbuser := range dbusers {
		usr, t := toCoreUserWithTotal(dbuser)
		users[i] = usr
		total = t
	}
	return users, total
}
