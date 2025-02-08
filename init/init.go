package init

import (
	"context"
	"fmt"

	oms_kafka "oms-service/kafka"
	"time"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func Initialize(ctx context.Context) {
	initializeLog(ctx)
	initializeDB(ctx)
	// initializeRedis(ctx)
	oms_kafka.InitializeKafkaProducer()
	go oms_kafka.InitializeKafkaConsumer(ctx)
	initialiseSQSProducer(ctx)
	initialiseSQSConsumer(ctx)
}

// Initialize logging
func initializeLog(ctx context.Context) {
	err := log.InitializeLogger(
		log.Formatter(config.GetString(ctx, "log.format")),
		log.Level(config.GetString(ctx, "log.level")),
	)
	if err != nil {
		log.WithError(err).Panic("unable to initialise log")
	}
}

func initializeDB(c context.Context) {
	fmt.Println("Connecting to mongo...")
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(getDatabaseUri())

	var err error
	DB, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}

	err = DB.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Failed to ping MongoDB:", err)
		return
	}

	fmt.Println("Successfully connected to MongoDB!")
}

func getDatabaseUri() string {
	return "mongodb://localhost:27017"
}

// Initialize Redis
// func initializeRedis(ctx context.Context) {
// 	r := oredis.NewClient(&oredis.Config{
// 		ClusterMode: config.GetBool(ctx, "redis.clusterMode"),
// 		Hosts:       config.GetStringSlice(ctx, "redis.hosts"),
// 		DB:          config.GetUint(ctx, "redis.db"),
// 	})
// 	log.InfofWithContext(ctx, "Initialized Redis Client")
// 	redis.SetClient(r)
// }
