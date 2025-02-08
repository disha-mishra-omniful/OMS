package database

import "go.mongodb.org/mongo-driver/mongo"

type Connect struct {
	*mongo.Client
}

var DB *Connect

func GetClient() *Connect {
	return DB
}
func SetClient(client *mongo.Client) {
	DB = &Connect{client}
}
