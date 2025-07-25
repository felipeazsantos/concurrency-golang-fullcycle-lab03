package bidcontroller

import (
	"context"
	"net/http"

	resterr "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/configuration/rest_err"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (b *BidController) FindBidByAuctionId(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		errRest := resterr.NewBadRequestError("invalid fields", resterr.Causes{
			Field:   "auctionId",
			Message: "invalid UUID value",
		})

		c.JSON(errRest.Code, errRest)
		return
	}

	bidOutputList, err := b.bidUseCase.FindBidByAuctionId(context.Background(), auctionId)
	if err != nil {
		restErr := resterr.ConvertError(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, bidOutputList)
}

func (b *BidController) FindWinningBidByAuctionId(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		restErr := resterr.NewBadRequestError("invalid UUID", resterr.Causes{
			Field:   "auctionId",
			Message: "invalid UUID",
		})
		c.JSON(restErr.Code, restErr)
		return
	}

	bidWinning, err := b.bidUseCase.FindWinningBidByAuctionId(context.Background(), auctionId)
	if err != nil {
		restErr := resterr.ConvertError(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, bidWinning)
}
