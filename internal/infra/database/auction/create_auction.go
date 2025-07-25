package auction

import (
	"context"
	"os"
	"time"

	"github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/configuration/logger"
	auctionentity "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/entity/auction_entity"
	internalerror "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollectionInterface interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
}

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
	Collection CollectionInterface
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func NewAuctionRepositoryWithCollection(collection CollectionInterface) *AuctionRepository {
	return &AuctionRepository{
		Collection: collection,
	}
}

func (ar *AuctionRepository) CreateAuction(ctx context.Context, auctionEntity *auctionentity.Auction) *internalerror.InternalError {
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

	go func() {
		<-time.After(getAuctionInterval())
		err := ar.CloseAuction(context.Background(), auctionEntityMongo.Id)
		if err != nil {
			logger.Error("error trying to update auction status to completed", err)
			return
		}
	}()

	return nil
}

func (ar *AuctionRepository) CloseAuction(ctx context.Context, auctionId string) *internalerror.InternalError {
	update := bson.M{"$set": bson.M{"status": auctionentity.Completed}}
	filter := bson.M{"_id": auctionId}

	_, err := ar.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return internalerror.NewInternalServerError("error trying to update auction")
	}

	return nil
}

func getAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return 5 * time.Minute
	}
	return duration
}
