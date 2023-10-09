package token

import (
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{
		Id:        id,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

func (payload *Payload) Valid() error {
	if time.Now().Unix() > payload.ExpiredAt.Unix() {
		return ErrExpiredToken
	}

	return nil
}
