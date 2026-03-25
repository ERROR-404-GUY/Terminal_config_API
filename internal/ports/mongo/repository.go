package mongo

import (
	"context"
	"terminal_config/internal/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type terminalRepository struct {
	collection *mongo.Collection
}

func NewTerminalConfigRepository(db *mongo.Database) domain.TerminalRepository {
	return &terminalRepository{
		collection: db.Collection("terminal_config"),
	}
}

func (r *terminalRepository) CreateTerminal(ctx context.Context, terminal *domain.Terminal) error {
	_, err := r.collection.InsertOne(ctx, terminal)
	return err
}

func (r *terminalRepository) GetTerminalByTID(ctx context.Context, tid string) (*domain.Terminal, error) {
	var terminal domain.Terminal
	err := r.collection.FindOne(ctx, bson.M{"tid": tid}).Decode(&terminal)
	if err != nil {
		return nil, err
	}
	return &terminal, nil
}

func (r *terminalRepository) UpdateTerminal(ctx context.Context, terminal *domain.Terminal) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"tid": terminal.TID},
		bson.M{"$set": terminal},
	)
	return err
}

func (r *terminalRepository) DeleteTerminal(ctx context.Context, tid string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"tid": tid})
	return err
}

func (r *terminalRepository) ListTerminals(ctx context.Context) ([]*domain.Terminal, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var terminals []*domain.Terminal
	if err = cursor.All(ctx, &terminals); err != nil {
		return nil, err
	}
	return terminals, nil
}
