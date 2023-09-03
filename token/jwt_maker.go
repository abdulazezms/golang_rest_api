package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

)

const (
	minKeySize = 32
)

// JWTMaker is a json web token maker that implements the maker interface
type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minKeySize {
		return nil, fmt.Errorf("Invalid key size: Must be at least %d characters", minKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	fmt.Println(payload)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		//check if the hashing algorithm matches to what we use, return the key.
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, payload, keyFunc)
	if err != nil {
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
