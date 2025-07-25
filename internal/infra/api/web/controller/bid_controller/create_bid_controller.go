package bidcontroller

import (
	"context"
	"net/http"

	resterr "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/configuration/rest_err"
	"github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/infra/api/web/validation"
	bidusecase "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/usecase/bid_usecase"
	"github.com/gin-gonic/gin"
)

type BidController struct {
	bidUseCase bidusecase.BidUseCaseInterface
}

func NewbidController(bidUseCase bidusecase.BidUseCaseInterface) *BidController {
	return &BidController{
		bidUseCase: bidUseCase,
	}
}

func (a *BidController) Createbid(c *gin.Context) {
	var bidInputDTO bidusecase.BidInputDTO

	if err := c.ShouldBindJSON(&bidInputDTO); err != nil {
		restErr := validation.ValidateErr(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	if err := a.bidUseCase.CreateBid(context.Background(), bidInputDTO); err != nil {
		restErr := resterr.ConvertError(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	c.Status(http.StatusCreated)
}
