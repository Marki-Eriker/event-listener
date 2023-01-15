package token

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/marki-eriker/event-listener/entity/user"
	"time"
)

type JWT struct {
	secret string
	ttl    time.Duration
}

func NewJWT(opt *Options) *JWT {
	return &JWT{secret: opt.Secret, ttl: opt.Ttl}
}

type Claims struct {
	jwt.RegisteredClaims
	Role user.Role `json:"role"`
}

func (j *JWT) Generate(role user.Role) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ttl)),
		},
		Role: role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secret))
}

func (j *JWT) Parse(rawToken string) (*jwt.Token, error) {
	jwtToken, err := jwt.ParseWithClaims(rawToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	return jwtToken, nil
}
