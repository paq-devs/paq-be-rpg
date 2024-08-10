package config

import (
	"log"

	"github.com/paq-devs/paq-be-rpg/api/repository"
	"github.com/paq-devs/paq-be-rpg/internal/lobby"
)

type Module struct {
	LobbyService *lobby.LobbyService
}

var module = Module{}

func Init() {
	mongoCfg := MongoConfig{ // TO-DO: Environment variables
		URI:            "mongodb://localhost:27017",
		DatabaseName:   "paq_rpg",
		CollectionName: "lobby",
	}

	db, _, err := ConnectMongoDB(mongoCfg)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewMongoLobbyRepository(db, mongoCfg.CollectionName)
	module.LobbyService = lobby.NewLobbyService(repo)
}

func GetModule() *Module {
	return &module
}
