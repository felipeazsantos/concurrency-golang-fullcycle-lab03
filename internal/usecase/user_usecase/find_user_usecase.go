package userusecase

import (
	"context"

	userentity "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/entity/user_entity"
	internalerror "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/internal_error"
)

type UserUseCase struct {
	UserRepository userentity.UserRepositoryInterface
}

type UserOutputDTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserUseCaseInterface interface {
	FindUserById(ctx context.Context, userId string) (*UserOutputDTO, *internalerror.InternalError)
}

func NewUserUseCase(userRepository userentity.UserRepositoryInterface) UserUseCaseInterface {
	return &UserUseCase{
		UserRepository: userRepository,
	}
}

func (u *UserUseCase) FindUserById(ctx context.Context, userId string) (*UserOutputDTO, *internalerror.InternalError) {

	user, err := u.UserRepository.FindUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &UserOutputDTO{
		Id:   user.Id,
		Name: user.Name,
	}, nil
}
