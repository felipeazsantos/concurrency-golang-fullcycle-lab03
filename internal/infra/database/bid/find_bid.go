package bid

import (
	"context"
	"fmt"
	"time"

	"github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/configuration/logger"
	bidentity "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/entity/bid_entity"
	internalerror "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (br *BidRepository) FindBidByAuctionId(
	ctx context.Context,
	auctionId string,
) ([]bidentity.Bid, *internalerror.InternalError) {
	filter := bson.M{"auctionId": auctionId}

	cursor, err := br.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error(fmt.Sprintf("error trying to find bids by auction id = %s", auctionId), err)
		return nil, internalerror.NewInternalServerError(fmt.Sprintf("error trying to find bids by auction id = %s", auctionId))
	}

	var bidsEntityMongo []BidEntityMongo
	if err := cursor.All(ctx, &bidsEntityMongo); err != nil {
		logger.Error(fmt.Sprintf("error trying to find bids by auction id = %s", auctionId), err)
		return nil, internalerror.NewInternalServerError(fmt.Sprintf("error trying to find bids by auction id = %s", auctionId))
	}

	var bidsEntity []bidentity.Bid
	for _, bid := range bidsEntityMongo {
		bidsEntity = append(bidsEntity, bidentity.Bid{
			Id:        bid.Id,
			UserId:    bid.UserId,
			AuctionId: bid.AuctionId,
			Amount:    bid.Amount,
			Timestamp: time.Unix(bid.Timestamp, 0),
		})
	}

	return bidsEntity, nil
}

func (br *BidRepository) FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*bidentity.Bid, *internalerror.InternalError) {
	filter := bson.M{"auctionId": auctionId}

	opts := options.FindOne().SetSort(bson.D{{Key: "amount", Value: -1}})
	var bidEntityMongo BidEntityMongo
	if err := br.Collection.FindOne(ctx, filter, opts).Decode(&bidEntityMongo); err != nil {

		logger.Error(fmt.Sprintf("error trying to find winnig bid by auction id = %s", auctionId), err)
		return nil, internalerror.NewInternalServerError(fmt.Sprintf("error trying to find winnig bid by auction id = %s", auctionId))
	}

	return &bidentity.Bid{
		Id:        bidEntityMongo.Id,
		UserId:    bidEntityMongo.UserId,
		AuctionId: bidEntityMongo.AuctionId,
		Amount:    bidEntityMongo.Amount,
		Timestamp: time.Unix(bidEntityMongo.Timestamp, 0),
	}, nil
}
