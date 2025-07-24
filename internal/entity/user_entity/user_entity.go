package userentity

import (
	"context"

	internalerror "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/internal_error"
)

type User struct {
	Id   string
	Name string
}

type UserRepositoryInterface interface {
	FindUserById(ctx context.Context, userId string) (*User, *internalerror.InternalError)
}
