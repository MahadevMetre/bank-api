package app

import (
	"bankapi/config"
	"bankapi/constants"
	"bankapi/middleware"
	"bankapi/router"
	"bankapi/rpc"
	"bankapi/stores"
	"context"
	"fmt"
	"log"

	"bitbucket.org/paydoh/paydoh-commons/amazon"
	"bitbucket.org/paydoh/paydoh-commons/database"
	"bitbucket.org/paydoh/paydoh-commons/services"
	"bitbucket.org/paydoh/paydoh-commons/settings"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SetupTestRouter() (*gin.Engine, func()) {
	appCtx, cncl := context.WithCancel(context.Background())

	settings.LoadEnvFile()

	awsKeyId := constants.AWSAccessKeyID
	awsSecretKey := constants.AWSSecretAccessKey
	awsRegion := constants.AWSRegion

	awsInstance, err := amazon.NewAws(
		awsRegion,
		awsKeyId,
		awsSecretKey,
		constants.AWSBucketName,
		constants.AWSCloudFrontURL,
	)
	if err != nil {
		log.Fatal(err)
	}

	// logger
	loggerSrv, err := services.NewLoggerService(&services.Config{
		AwsRegion:          awsRegion,
		AwsAccessKeyID:     awsKeyId,
		AwsSecretAccessKey: awsSecretKey,
		EnableRequestLog:   true,
		EnableResponseLog:  true,
		EnableS3Upload:     false,
	})
	if err != nil {
		log.Println(err)
		panic(err)
	}

	// redis
	memoryConfig := database.InMemoryConfig{
		NetworkType: constants.RedisNetworkType,
		Address:     constants.RedisURL,
		Username:    constants.RedisUserName,
		Password:    constants.RedisPassword,
	}
	memory, err := database.NewInMemory(appCtx, memoryConfig)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	// mongodb
	client, err := database.ConnectToMongoDB(appCtx, constants.MongodbURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	dbdoc, err := database.NewDocument(client)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// postgres
	db := config.InitDB()

	// goose migrations
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Failed to set dialect: %v", err)
	}

	migrationDir := "../../migrations"
	if err := goose.Up(db, migrationDir); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	fmt.Println("Migrations applied successfully!")

	grpcOptions := []grpc.DialOption{}
	grpcOptions = append(grpcOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(constants.PaymentServiceURL, grpcOptions...)
	if err != nil {
		panic(err)
	}

	paymentClient := rpc.NewPaymentServiceClient(conn)

	s := stores.NewStores(loggerSrv, db, memory, dbdoc, appCtx, paymentClient)

	app := gin.Default()

	app.Use(cors.Default())
	app.Use(s.BindStore(awsInstance))
	app.Use(middleware.LoggerMiddleware(loggerSrv))
	router.Endpoints(app)

	s.StartPeriodicTask()

	// Cleanup function
	cleanup := func() {
		cncl()
		conn.Close()
		db.Close()
		client.Disconnect(appCtx)
		memory.Close()
		fmt.Println("Cleanup completed")
	}

	return app, cleanup
}
