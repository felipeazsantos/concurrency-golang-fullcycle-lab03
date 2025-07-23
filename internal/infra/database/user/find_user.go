package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/configuration/logger"
	userentity "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/entity/user_entity"
	internalerror "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserEntityMongo struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection: database.Collection("users"),
	}
}

func (ur *UserRepository) FindUserById(ctx context.Context, userId string) (*userentity.User, *internalerror.InternalError) {
	filter := bson.M{"_id": userId}

	var userEntityMongo UserEntityMongo
	err := ur.Collection.FindOne(ctx, filter).Decode(&userEntityMongo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error(fmt.Sprintf("User not found with this id = %s ", userId), err)
			return nil, internalerror.NewNotFoundError(fmt.Sprintf("User not found with this id = %s ", userId))
		}

		logger.Error("Error trying to find user by userId", err)
		return nil, internalerror.NewInternalServerError("Error trying to find user by userId")
	}

	userEntity := &userentity.User{
		Id:   userEntityMongo.Id,
		Name: userEntityMongo.Name,
	}

	return userEntity, nil
}
