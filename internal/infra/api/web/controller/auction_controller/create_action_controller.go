package auctioncontroller

import (
	"context"
	"net/http"

	resterr "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/configuration/rest_err"
	"github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/infra/api/web/validation"
	auctionusecase "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/usecase/auction_usecase"
	"github.com/gin-gonic/gin"
)

type AuctionController struct {
	auctionUseCase auctionusecase.AuctionUseCaseInterface
}

func NewAuctionController(auctionUseCase auctionusecase.AuctionUseCaseInterface) *AuctionController {
	return &AuctionController{
		auctionUseCase: auctionUseCase,
	}
}

func (a *AuctionController) CreateAuction(c *gin.Context) {
	var auctionInputDTO auctionusecase.AuctionInputDTO

	if err := c.ShouldBindJSON(&auctionInputDTO); err != nil {
		restErr := validation.ValidateErr(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	if err := a.auctionUseCase.CreateAuction(context.Background(), auctionInputDTO); err != nil {
		restErr := resterr.ConvertError(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	c.Status(http.StatusCreated)
}
