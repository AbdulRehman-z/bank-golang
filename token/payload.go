package token

import (
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  int64     `json:"issued_at"`
	ExpiredAt int64     `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{
		Id:        id,
		Username:  username,
		IssuedAt:  time.Now().Unix(),
		ExpiredAt: time.Now().Add(duration).Unix(),
	}, nil
}

func (payload *Payload) Valid() error {
	if time.Now().Unix() > payload.ExpiredAt {
		return ErrExpiredToken
	}

	return nil
}
