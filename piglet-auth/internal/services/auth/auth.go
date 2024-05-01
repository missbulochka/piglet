package auth

import (
	"context"
	"log/slog"
)

type Auth struct {
	log *slog.Logger
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		username string,
		email string,
		passHash []byte,
	) (uid int64, err error)
	UpdateCurrUser(
		ctx context.Context,
		username string,
		email string,
		oldPassHash []byte,
		passHash []byte,
	) (err error)
}

type UserProvider interface {
	User(ctx context.Context, username string) (uid int64, err error)
}
