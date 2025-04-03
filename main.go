package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bankapi/config"
	"bankapi/constants"
	"bankapi/middleware"
	"bankapi/router"
	"bankapi/rpc"
	"bankapi/stores"

	"bitbucket.org/paydoh/paydoh-commons/amazon"
	"bitbucket.org/paydoh/paydoh-commons/database"
	"bitbucket.org/paydoh/paydoh-commons/pkg/worker"
	"bitbucket.org/paydoh/paydoh-commons/services"
	"bitbucket.org/paydoh/paydoh-commons/settings"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	responseMiddleware "bitbucket.org/paydoh/paydoh-commons/middleware"
)

var (
	appCtx context.Context
	cncl   context.CancelFunc
)

// @title Paydoh bank api
// @version 1.0
// @contact.name Paydoh
// @contact.email it@paydoh.money
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:4100
// @schemes http https
// @basePath /
func main() {
	appCtx, cncl = context.WithCancel(context.Background())

	defer cncl()

	settings.LoadEnvFile()

	app := gin.Default()

	// Register pprof handlers
	pprof.Register(app)

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
		panic(err)
	}

	loggerSrv, err := services.NewLoggerService(&services.Config{
		AwsRegion:          awsRegion,
		AwsAccessKeyID:     awsKeyId,
		AwsSecretAccessKey: awsSecretKey,
		EnableRequestLog:   true,
		EnableResponseLog:  true,
		EnableS3Upload:     false,
		LogLevel:           constants.LOG_LEVEL,
	})
	if err != nil {
		log.Println(err)
		panic(err)
	}

	memory, err := config.InitRedis(appCtx)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	client, err := database.ConnectToMongoDB(appCtx, constants.MongodbURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	dbdoc, err := database.NewDocument(client)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	db := config.InitDB()

	// goose migrations
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Failed to set dialect: %v", err)
	}

	migrationDir := "./migrations"
	if err := goose.Up(db, migrationDir); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	fmt.Println("Migrations applied successfully!")

	defer db.Close()

	grpcOptions := []grpc.DialOption{}
	grpcOptions = append(grpcOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(constants.PaymentServiceURL, grpcOptions...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	paymentClient := rpc.NewPaymentServiceClient(conn)

	s := stores.NewStores(loggerSrv, db, memory, dbdoc, appCtx, paymentClient)

	app.Static("/static", "./templates/static")

	xssMiddleware := responseMiddleware.NewXSSMiddlewareConfig()
	app.Use(xssMiddleware.XssMiddlewareProtect())
	app.Use(gzip.Gzip(gzip.DefaultCompression))
	app.Use(s.BindStore(awsInstance))
	app.Use(middleware.LoggerMiddleware(loggerSrv))
	app.Use(responseMiddleware.ResponseHeaderMiddleware())
	router.Endpoints(app)

	s.StartPeriodicTask()

	defer loggerSrv.Stop()

	// scheduler, err := ec2_scheduler.NewScheduler(ec2_scheduler.AWSCredentials{
	// 	AccessKeyID:     awsKeyId,
	// 	SecretAccessKey: awsSecretKey,
	// 	Region:          awsRegion,
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to create scheduler: %v", err)
	// }

	// go func() {
	// 	if err := scheduler.Start(); err != nil {
	// 		log.Fatalf("Failed to start scheduler: %v", err)
	// 	}
	// }()
	// defer scheduler.Stop()

	serverConfig := config.GetServerConfig()
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", constants.GolangPort),
		Handler:        app.Handler(),
		ReadTimeout:    serverConfig.ReadTimeout,
		WriteTimeout:   serverConfig.WriteTimeout,
		IdleTimeout:    serverConfig.IdleTimeout,
		MaxHeaderBytes: serverConfig.MaxHeaderBytes,
	}
	app.LoadHTMLGlob("templates/*.html")

	go func() {
		fmt.Printf("Server is running on http://localhost:%d\n", constants.GolangPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	asynq := worker.NewAsynqServer(constants.RedisURL, constants.RedisUserName, constants.RedisPassword, constants.RedisDB, 10)
	defer asynq.Stop()

	go func() {
		if err := asynq.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	asynq.HandleFunc(constants.AuditLogType, s.AuditLogService.AuditLogHandler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 2 seconds.")
	}
	log.Println("Server exiting")
}
