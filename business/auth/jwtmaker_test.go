package auth

import (
	"github/islamghany/blog/business/core/user"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const (
	secretKey = "secretssdfkjsdsdfjhsdfsdjkfsdfhsdkfsdkhfsdfsdfsdfsdfkjsdfdfhyww"
)

// go test -v ./business/auth -run TestJWTMaker

func createToken(t *testing.T, maker Maker) (string, *Payload) {

	id := uuid.New()
	roles := []user.Role{user.RoleAdmin}
	duration := time.Second * 10
	tk, p, err := maker.Sign(id, roles, 1, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tk)
	require.NotEmpty(t, p)
	require.Equal(t, id.String(), p.ID.String())
	require.Equal(t, roles, p.Roles)
	// t.Errorf("IssuedAt: %+v\n", p)

	return tk, p
}
func TestTokenSigning(t *testing.T) {
	// log := NewLogger()
	// maker := NewAuth(secretKey)
	// tk, payload := createToken(t, maker)
	// t.Logf("Token: %v\n", tk)
	// t.Logf("Payload: %v\n", payload)
	// payload, err := maker.Verify(tk)
	// // time.Sleep(5 * time.Second)
	// require.NoError(t, err)
	// require.NotEmpty(t, payload)
	// require.Equal(t, payload.ID.String(), payload.ID.String())
	// require.Equal(t, payload.Roles, payload.Roles)
	// require.WithinDuration(t, payload.IssuedAt, time.Now(), time.Second)
	// require.WithinDuration(t, payload.ExpiredAt, time.Now().Add(10*time.Second), time.Second)
	// require.Equal(t, payload.Version, payload.Version)

}
