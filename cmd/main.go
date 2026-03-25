package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"terminal_config/internal/domain"
	ports "terminal_config/internal/ports/http"
	configmongo "terminal_config/internal/ports/mongo"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func main() {
	mongoClient, err := setupMongoClient()
	if err != nil {
		panic(err)
	}

	terminalConfigRepo := configmongo.NewTerminalConfigRepository(mongoClient.Database("ExperienceCluster0"))
	terminalConfigService := domain.NewTerminalConfigService(terminalConfigRepo)

	muxRouter := mux.NewRouter()

	handler := ports.NewHandlers(muxRouter, terminalConfigService)
	handler.SetupRoutes()

	log.Fatal(http.ListenAndServe(":8081", muxRouter))

	// Now you can use terminalConfigService to handle requests

}

func setupMongoClient() (*mongo.Client, error) {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://admin:kMuKWnuOv9rMmYB8@experiencecluster0.hpzvuqf.mongodb.net/?appName=ExperienceCluster0").SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return client, nil
}
