package auction

import (
	"context"

	"github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/configuration/logger"
	auctionentity "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/entity/auction_entity"
	internalerror "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/internal_error"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                         `bson:"_id"`
	ProductName string                         `bson:"product_name"`
	Category    string                         `bson:"category"`
	Description string                         `bson:"description"`
	Condition   auctionentity.ProductCondition `bson:"condition"`
	Status      auctionentity.AuctionStatus    `bson:"status"`
	Timestamp   int64                          `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(ctx context.Context, auctionEntity auctionentity.Auction) *internalerror.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}

	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internalerror.NewInternalServerError("Error trying to insert auction ")
	}

	return nil
}
