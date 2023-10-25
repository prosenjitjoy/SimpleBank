package token

import (
	"encoding/json"
	"errors"
	"time"

	"aidanwoods.dev/go-paseto"
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	secretKey paseto.V4SymmetricKey
}

// NewPASETOMaker creates a new PasetoMaker
func NewPASETOMaker(secretKey string) (Maker, error) {
	if len(secretKey) != 32 {
		return nil, errors.New("invalid key size: must be exactly 32 characters")
	}

	symmetricKey, err := paseto.V4SymmetricKeyFromBytes([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &PasetoMaker{
		secretKey: symmetricKey,
	}, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
	}

	token := paseto.NewToken()
	token.Set("id", payload.ID)
	token.Set("username", payload.Username)
	token.SetIssuedAt(payload.IssuedAt)
	token.SetExpiration(payload.ExpiredAt)

	return token.V4Encrypt(maker.secretKey, nil), payload, nil
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	parser := paseto.NewParser()
	newToken, err := parser.ParseV4Local(maker.secretKey, token, nil)
	if err != nil {
		if err.Error() == "this token has expired" {
			return nil, ErrExpiredToken
		}
		return nil, err
	}

	payload := &Payload{}
	err = json.Unmarshal(newToken.ClaimsJSON(), payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
