package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	URI            string
	DatabaseName   string
	CollectionName string
}

func ConnectMongoDB(cfg MongoConfig) (*mongo.Database, *mongo.Collection, error) {
	clientOptions := options.Client().ApplyURI(cfg.URI)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao conectar ao MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao verificar a conexão com MongoDB: %v", err)
	}

	log.Println("Conectado ao MongoDB!")

	db := client.Database(cfg.DatabaseName)
	collection := db.Collection(cfg.CollectionName)

	return db, collection, nil
}
