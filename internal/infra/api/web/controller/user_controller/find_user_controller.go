package usercontroller

import (
	"context"
	"net/http"

	resterr "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/configuration/rest_err"
	userusecase "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/usecase/user_usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userController struct {
	userUseCase userusecase.UserUseCaseInterface
}

func NewUserController(userUseCase userusecase.UserUseCaseInterface) *userController {
	return &userController{
		userUseCase: userUseCase,
	}
}

func (u *userController) FindUserById(c *gin.Context) {
	userId := c.Param("userId")

	if err := uuid.Validate(userId); err != nil {
		errRest := resterr.NewBadRequestError("invalid fields", resterr.Causes{
			Field:   "userId",
			Message: "invalid UUID value",
		})

		c.JSON(errRest.Code, errRest)
		return
	}

	userData, err := u.userUseCase.FindUserById(context.Background(), userId)
	if err != nil {
		restErr := resterr.ConvertError(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, userData)
}
