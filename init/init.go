package appinit

import (
	"awesomeProject5/OMS/database"
	"awesomeProject5/OMS/kafkaa"
	"awesomeProject5/OMS/orders/listeners"
	//
	"awesomeProject5/OMS/orders/services"
	"awesomeProject5/OMS/redis"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/omniful/go_commons/kafka"
	"log"
	"os"

	//"github.com/joho/godotenv"
	"github.com/omniful/go_commons/config"
	goredis "github.com/omniful/go_commons/redis"
	"github.com/omniful/go_commons/sqs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"log"
	//"os"
	"time"
)

//contextKey string

// const DBKey contextKey = "mongoDB"

// var DB *mongo.Client

func Initialize(ctx context.Context) {
	InitializeRedis(ctx)
	InitializeDB(ctx)
	//InitializeKafka(ctx)
	InitializeSQS(ctx)

}

// err:=godotenv.load()
//
//	if err != nil {
//		log.Fatal("Error loading .env file")
//	}
func InitializeDB(ctx context.Context) {
	fmt.Println("Connecting to mongo...")
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(config.GetString(ctx, "mongo.string"))

	var err error
	Db, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}
	err = Db.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Failed to ping MongoDB:", err)
		return
	}

	fmt.Println("Successfully connected to MongoDB!")

	database.SetClient(Db)
}

// Initialize Kafka Producer
func initializeKafkaProducer(ctx context.Context) {
	kafkaBrokers := config.GetStringSlice(ctx, "onlineKafka.brokers")
	kafkaClientID := config.GetString(ctx, "onlineKafka.clientId")
	kafkaVersion := config.GetString(ctx, "onlineKafka.version")

	producer := kafka.NewProducer(
		kafka.WithBrokers(kafkaBrokers),
		kafka.WithClientID(kafkaClientID),
		kafka.WithKafkaVersion(kafkaVersion),
	)
	log.Printf("Initialized Kafka Producer")
	//kafka.Set(producer)
	kafkaa.Set(producer)
	//go listeners.InitializeKafkaConsumer(ctx)
	go listeners.InitializeKafkaConsumer(ctx)

}
func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("unable to load")
	}

}
func InitializeSQS(ctx context.Context) {
	SQSconfig := sqs.GetSQSConfig(ctx, false, "order", "eu-north-1", os.Getenv("AWS_ACCOUNT"), "")
	queue_url, err := sqs.GetUrl(ctx, SQSconfig, "sqsQueue")
	if err != nil {
		log.Fatal("cant get url")
	}
	// log.Printf("Successfully initialized SQS. Queue URL: %s", *queue_url)
	Queue_instance, err := sqs.NewStandardQueue(ctx, "sqsQueue", SQSconfig)
	if err != nil {
		log.Fatal("cant create queue instance")
	}
	fmt.Println(queue_url)
	fmt.Println(Queue_instance)
	// fmt.Println(Queue_instance,*Queue_instance.Url)
	services.SetProducer(ctx, Queue_instance)
	go listeners.StartConsume(*queue_url, ctx)

}
func InitializeRedis(ctx context.Context) {
	redis_client := goredis.NewClient(&goredis.Config{
		ClusterMode: config.GetBool(ctx, "redis.clusterMode"),
		Hosts:       []string{config.GetString(ctx, "redis.hosts")},
		DB:          config.GetUint(ctx, "redis.postgres"),
	})
	fmt.Println("Initialized redis client")

	redis.SetClient(redis_client)

}
