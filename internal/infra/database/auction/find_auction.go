package auction

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/configuration/logger"
	auctionentity "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/entity/auction_entity"
	internalerror "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (ar *AuctionRepository) FindAuctionById(ctx context.Context, id string) (*auctionentity.Auction, *internalerror.InternalError) {
	filter := bson.M{"_id": id}

	var auctionEntityMongo AuctionEntityMongo
	err := ar.Collection.FindOne(ctx, filter).Decode(&auctionEntityMongo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error(fmt.Sprintf("Error trying find user with this id = %s", id), err)
			return nil, internalerror.NewNotFoundError(fmt.Sprintf("Error trying find auction with this id = %s", id))
		}

		logger.Error("Error trying to find user by id", err)
		return nil, internalerror.NewInternalServerError("Error trying to find auction by id")
	}

	auctionEntity := &auctionentity.Auction{
		Id:          auctionEntityMongo.Id,
		ProductName: auctionEntityMongo.ProductName,
		Category:    auctionEntityMongo.Category,
		Description: auctionEntityMongo.Description,
		Condition:   auctionEntityMongo.Condition,
		Status:      auctionEntityMongo.Status,
		Timestamp:  time.Unix(auctionEntityMongo.Timestamp, 0),
	}

	return auctionEntity, nil
}

func (ar *AuctionRepository) FindAuctions(
	ctx context.Context, 
	status auctionentity.AuctionStatus,
	category string,
	productName string,
) ([]auctionentity.Auction, *internalerror.InternalError) {
	filter := bson.M{}

	if status != 0 {
		filter["status"] = status
	}

	if category != "" {
		filter["category"] = category
	}

	if productName != "" {
		filter["productName"] = primitive.Regex{
			Pattern: productName,
			Options: "i",
		}
	}

	 cursor, err := ar.Collection.Find(ctx, filter)
	 if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error("Error trying to find auctions not found", err)
			return nil, internalerror.NewNotFoundError("Error trying to find auctions not found")
		}

		logger.Error("Error trying find auctions", err)
		return nil, internalerror.NewInternalServerError("Error trying to find auctions")
	 }

	 defer cursor.Close(ctx)

	 var auctionsEntityMongo []AuctionEntityMongo
	 if err := cursor.All(ctx, &auctionsEntityMongo); err != nil {
		logger.Error("Error trying to find auctions", err)
		return nil, internalerror.NewInternalServerError("Error trying to find auctions")
	 }

	 var auctionsEntity []auctionentity.Auction
	 for _, auction := range auctionsEntityMongo {

		auctionsEntity = append(auctionsEntity, auctionentity.Auction{
			Id:          auction.Id,
			ProductName: auction.ProductName,
			Category:    auction.Category,
			Description: auction.Description,
			Condition:   auction.Condition,
			Status:      auction.Status,
			Timestamp:   time.Unix(auction.Timestamp, 0),
		})
	 }

	 return auctionsEntity, nil
} 
