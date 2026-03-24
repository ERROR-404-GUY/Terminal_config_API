package mongo

import "go.mongodb.org/mongo-driver/v2/mongo"

type TerminalConfigRepository interface {
	// Define methods for interacting with the terminal configuration data
}

type terminalConfigRepository struct {
	collection *mongo.Collection
}

func NewTerminalConfigRepository(db *mongo.Database) TerminalConfigRepository {
	return &terminalConfigRepository{
		collection: db.Collection("terminal_config"),
	}
}
