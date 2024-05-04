package user

import "fmt"

var (
	RoleAdmin = Role{"admin"}
	RoleUser  = Role{"user"}
)

// set of known roles

var roles = map[string]Role{
	RoleAdmin.name: RoleAdmin,
	RoleUser.name:  RoleUser,
}

type Role struct {
	name string
}

func ParseRole(value string) (Role, error) {
	role, ok := roles[value]
	if !ok {
		return Role{}, fmt.Errorf("'%s' is not a valid role", value)
	}
	return role, nil
}

// MustParseRole parses the string value and returns a role if one exists. If
// an error occurs the function panics.
func MustParseRole(value string) Role {
	role, err := ParseRole(value)
	if err != nil {
		panic(err)
	}

	return role
}

// Name returns the name of the role.
func (r Role) Name() string {
	return r.name
}

func (r Role) String() string {
	return r.name
}

func (r Role) MarshalText() ([]byte, error) {
	return []byte(r.name), nil
}

func (r *Role) UnmarshalText(data []byte) error {
	role, err := ParseRole(string(data))
	if err != nil {
		return err
	}

	r.name = role.name

	return nil
}
