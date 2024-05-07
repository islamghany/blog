package auth

import (
	"github/islamghany/blog/business/core/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// type JWTMaker struct {
// 	secretKey string
// }

// // NewJWTMaker creates a new JWTMaker
// // secretKey is the secret key used to sign tokens should be a complex random string
// func NewJWTMaker(secretKey string) Maker {
// 	return &JWTMaker{secretKey}
// }

func (a *Auth) Sign(id uuid.UUID, roles []user.Role, version int, duration time.Duration) (string, *Payload, error) {
	payload := NewPayload(id, roles, version, duration)
	rolesStr := make([]string, len(roles))
	for i, role := range roles {
		rolesStr[i] = role.String()
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      payload.ID,
		"roles":   rolesStr,
		"iss":     payload.IssuedAt.Unix(),
		"exp":     payload.ExpiredAt.Unix(),
		"version": payload.Version,
	})
	token, err := jwtToken.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", nil, err
	}
	return token, payload, nil
}

func (a *Auth) Verify(token string) (*Payload, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(a.secretKey), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	if !jwtToken.Valid {
		return nil, ErrInvalidToken
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	roles := []user.Role{}
	for _, role := range claims["roles"].([]interface{}) {
		roles = append(roles, user.MustParseRole(role.(string)))
	}
	return &Payload{
		ID:        uuid.MustParse(claims["id"].(string)),
		Roles:     roles,
		IssuedAt:  time.Unix(int64(claims["iss"].(float64)), 0),
		ExpiredAt: time.Unix(int64(claims["exp"].(float64)), 0),
		Version:   int(claims["version"].(float64)),
	}, nil
}
