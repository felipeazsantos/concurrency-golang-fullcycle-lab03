package auctioncontroller

import (
	"context"
	"net/http"
	"strconv"

	resterr "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/configuration/rest_err"
	auctionusecase "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/usecase/auction_usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (a *AuctionController) FindAuctionById(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		errRest := resterr.NewBadRequestError("invalid fields", resterr.Causes{
			Field:   "auctionId",
			Message: "invalid UUID value",
		})

		c.JSON(errRest.Code, errRest)
		return
	}

	auctionData, err := a.auctionUseCase.FindAuctionById(context.Background(), auctionId)
	if err != nil {
		restErr := resterr.ConvertError(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, auctionData)
}

func (a *AuctionController) FindAuctions(c *gin.Context) {
	status := c.Query("status")
	category := c.Query("category")
	productName := c.Query("productName")

	statusNumber, err := strconv.Atoi(status)
	if err != nil {
		restErr := resterr.NewBadRequestError("invalid field type", resterr.Causes{
			Field:   "status",
			Message: "invalid field type",
		})
		c.JSON(restErr.Code, restErr)
		return
	}

	auctions, internalErr := a.auctionUseCase.FindAuctions(context.Background(), auctionusecase.AuctionStatus(statusNumber), category, productName)
	if err != nil {
		restErr := resterr.ConvertError(internalErr)
		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, auctions)
}

func (a *AuctionController) FindWinningBidByAuctionId(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		restErr := resterr.NewBadRequestError("invalid UUID", resterr.Causes{
			Field:   "auctionId",
			Message: "invalid UUID",
		})
		c.JSON(restErr.Code, restErr)
		return
	}

	bidWinning, err := a.auctionUseCase.FindWinningBidByAuctionId(context.Background(), auctionId)
	if err != nil {
		restErr := resterr.ConvertError(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, bidWinning)
}
