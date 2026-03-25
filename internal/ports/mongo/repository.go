package mongo

import (
	"context"
	"terminal_config/internal/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type terminalRepository struct {
	collection *mongo.Collection
}

func NewTerminalConfigRepository(db *mongo.Database) domain.TerminalRepository {
	collection := db.Collection("terminal_config")
	// Enforce unique tid at the DB level
	_, _ = collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "tid", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	return &terminalRepository{
		collection: collection,
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

func (r *terminalRepository) UpdateRefundAllowed(ctx context.Context, tid string, allowed bool) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"tid": tid},
		bson.M{"$set": bson.M{"refund_allowed": allowed}},
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
