package repository

import (
	"context"
	"errors"

	lobby_ "github.com/paq-devs/paq-be-rpg/internal/lobby"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoLobbyRepository struct {
	collection *mongo.Collection
}

func NewMongoLobbyRepository(db *mongo.Database, collectionName string) *MongoLobbyRepository {
	return &MongoLobbyRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *MongoLobbyRepository) Save(ctx context.Context, lobby *lobby_.Lobby) error {
	_, err := r.collection.InsertOne(ctx, NewLobbyBson(lobby))
	return err
}

func (r *MongoLobbyRepository) FindByAccessCode(ctx context.Context, accessCode string) (*lobby_.Lobby, error) {
	var lobby LobbyBson

	filter := bson.M{"accessCode": accessCode}
	err := r.collection.FindOne(ctx, filter).Decode(&lobby)

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("lobby_not_found")
	}
	return lobby.ToLobby(), err
}

func (r *MongoLobbyRepository) Update(ctx context.Context, lobby *lobby_.Lobby) error {
	filter := bson.M{"_id": lobby.ID}

	update := bson.M{
		"$set": NewLobbyBson(lobby),
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}
