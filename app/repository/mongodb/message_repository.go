package mongodb

import (
	"context"
	"message-system/app/domain"
	"message-system/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepository struct {
	Client     *mongo.Client
	Database   string
	Collection string
}

func NewMessageRepository(config *config.MongoDBConfig) (*MessageRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(config.URI).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	return &MessageRepository{
		Client:     client,
		Database:   config.Database,
		Collection: config.Collection,
	}, nil
}

func (r *MessageRepository) GetUnsentMessages(ctx context.Context, limit int) ([]domain.Message, error) {
	db := r.Client.Database(r.Database)
	collection := db.Collection(r.Collection)

	filter := bson.M{"isSent": false}
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []domain.Message
	for cursor.Next(ctx) {
		var message domain.Message
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageRepository) GetSentMessages(ctx context.Context) ([]domain.Message, error) {
	db := r.Client.Database(r.Database)
	collection := db.Collection(r.Collection)

	filter := bson.M{"isSent": true}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []domain.Message
	for cursor.Next(ctx) {
		var message domain.Message
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageRepository) MarkAsSent(ctx context.Context, id int) error {
	db := r.Client.Database(r.Database)
	collection := db.Collection(r.Collection)

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"isSent": true}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
