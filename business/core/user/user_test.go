package user_test

import (
	"context"
	"github/islamghany/blog/business/core/user"
	"github/islamghany/blog/business/data/dbtest"
	"github/islamghany/blog/foundation/random"
	"runtime/debug"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestUser(t *testing.T) {
	t.Run("Create", testUserCreate)
}

// =============================================================================
// TestUserCreate
func testUserCreate(t *testing.T) {
	seed := func(ctx context.Context, usr *user.Core) ([]user.User, error) {
		return []user.User{}, nil
	}

	test := dbtest.NewTest(t)
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		test.TearDown()
	}()
	api := test.CoreAPIs.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Log("Go seeding ...")
	_, err := seed(ctx, api)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	t.Log("Create a user ...")
	nu := createRandomUser()
	usr, err := api.Create(ctx, nu)
	require.NoError(t, err)
	require.NotZero(t, usr.ID)
	require.NotZero(t, usr.CreatedAt)
	require.NotZero(t, usr.UpdatedAt)
	require.Equal(t, nu.Email, usr.Email)
	require.Equal(t, nu.Username, usr.Username)
	require.Equal(t, nu.FirstName, usr.FirstName)
	require.Equal(t, nu.LastName, usr.LastName)
	require.NotZero(t, usr.PasswordHashed)

}

func createRandomUser() user.NewUser {
	pass := random.RandomPassword()
	roles := make([]user.Role, 1)
	userRole, err := user.ParseRole("user")
	if err != nil {
		panic(err)
	}
	roles[0] = userRole
	usr := user.NewUser{
		Email:             random.RandomEmail(),
		Username:          random.RandomName(),
		Password:          pass,
		FirstName:         random.RandomName(),
		LastName:          random.RandomName(),
		Roles:             roles,
		ConfirmedPassword: pass,
	}

	return usr
}
