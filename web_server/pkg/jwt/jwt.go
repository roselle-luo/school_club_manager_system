package jwt

import (
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	UserID   uint   `json:"uid"`
	RoleCode string `json:"role"`
	jwtv5.RegisteredClaims
}

func Sign(secret string, uid uint, role string, expires int64) (string, error) {
	c := Claims{
		UserID:   uid,
		RoleCode: role,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(time.Duration(expires) * time.Second)),
			IssuedAt:  jwtv5.NewNumericDate(time.Now()),
		},
	}
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, c)
	return token.SignedString([]byte(secret))
}

func Parse(secret string, tokenStr string) (*Claims, error) {
	t, err := jwtv5.ParseWithClaims(tokenStr, &Claims{}, func(token *jwtv5.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := t.Claims.(*Claims); ok && t.Valid {
		return claims, nil
	}
	return nil, jwtv5.ErrTokenInvalidClaims
}
