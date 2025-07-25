package main

import (
	"context"
	"log"

	"github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/configuration/database/mongodb"
	"github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/configuration/logger"
	auctioncontroller "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/infra/api/web/controller/auction_controller"
	bidcontroller "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/infra/api/web/controller/bid_controller"
	usercontroller "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/infra/api/web/controller/user_controller"
	"github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/infra/database/auction"
	"github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/infra/database/bid"
	"github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/infra/database/user"
	auctionusecase "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/usecase/auction_usecase"
	bidusecase "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/usecase/bid_usecase"
	userusecase "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/internal/usecase/user_usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		if err := godotenv.Load("cmd/auction/.env"); err != nil {
			log.Fatal("Error trying to load env variables")
		}
	}

	ctx := context.Background()

	mongo, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	userController, auctionController, bidController := initDependencies(mongo)

	router := gin.Default()

	router.GET("/auctions/:auctionId", auctionController.FindAuctionById)
	router.GET("/auctions", auctionController.FindAuctions)
	router.POST("/auctions", auctionController.CreateAuction)
	router.GET("/auction/winner/:auctionId", auctionController.FindWinningBidByAuctionId)
	router.GET("/bid/:auctionId", bidController.FindBidByAuctionId)
	router.POST("/bid", bidController.Createbid)
	router.GET("/user/:userId", userController.FindUserById)

	logger.Info("starting server on port :8080")

	router.Run(":8080")
}

func initDependencies(database *mongo.Database) (
	userControler *usercontroller.UserController,
	auctionController *auctioncontroller.AuctionController,
	bidController *bidcontroller.BidController,
) {
	auctionRepository := auction.NewAuctionRepository(database)
	bidRepository := bid.NewBidRepository(database, auctionRepository)
	userRepository := user.NewUserRepository(database)

	auctionUseCase := auctionusecase.NewAuctionUseCase(auctionRepository, bidRepository)
	bidUseCase := bidusecase.NewBidUseCase(bidRepository)
	userUseCase := userusecase.NewUserUseCase(userRepository)

	userControler = usercontroller.NewUserController(userUseCase)
	auctionController = auctioncontroller.NewAuctionController(auctionUseCase)
	bidController = bidcontroller.NewbidController(bidUseCase)
	return
}
