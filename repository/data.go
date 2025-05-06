package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoRepository defines generic CRUD operations.
type MongoRepository[T any] interface {
	Insert(ctx context.Context, doc T) (*T, error)
	GetByID(ctx context.Context, id string) (*T, error)
	UpdateByID(ctx context.Context, id string, update T) error
	DeleteByID(ctx context.Context, id string) error
	FindAll(ctx context.Context, filter any) ([]T, error)
	CountDocuments(ctx context.Context, filter any) (int64, error)
}

// MongoRepositoryImpl provides a generic implementation for any model T.
type MongoRepositoryImpl[T any] struct {
	collection *mongo.Collection
}

func NewMongoRepository[T any](coll *mongo.Collection) MongoRepository[T] {
	return &MongoRepositoryImpl[T]{collection: coll}
}

func (r *MongoRepositoryImpl[T]) CountDocuments(ctx context.Context, filter any) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *MongoRepositoryImpl[T]) Insert(ctx context.Context, doc T) (*T, error) {
	_, err := r.collection.InsertOne(ctx, doc) // TODO : optimize - use response to validate data
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *MongoRepositoryImpl[T]) GetByID(ctx context.Context, id string) (*T, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}

	var result T
	err = r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *MongoRepositoryImpl[T]) UpdateByID(ctx context.Context, id string, update T) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (r *MongoRepositoryImpl[T]) DeleteByID(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (r *MongoRepositoryImpl[T]) FindAll(ctx context.Context, filter any) ([]T, error) {
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []T
	for cursor.Next(ctx) {
		var elem T
		if err := cursor.Decode(&elem); err != nil {
			return nil, err
		}
		results = append(results, elem)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
