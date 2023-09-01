package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

// JWTMaker is a json web token maker that implements the maker interface
type PasetoMaker struct {
	paseto    *paseto.V2
	secretKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) < chacha20poly1305.KeySize {
		return nil, fmt.Errorf("Invalid key size: Must be %d characters", chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto:    paseto.NewV2(),
		secretKey: []byte(symmetricKey),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil
	}
	return maker.paseto.Encrypt(maker.secretKey, payload, nil)
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.secretKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	
	if err != nil {
		return nil, ErrExpiredToken
	}
	return payload, nil
}
