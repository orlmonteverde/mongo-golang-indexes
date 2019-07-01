package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	// Name of the database.
	DBName          = "glottery"
	notesCollection = "notes"
	URI             = "mongodb://<user>:<password>@<host>/<name>"
)

type Note struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Title     string             `json:"title"`
	Body      string             `json:"body"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
}

var notes = []interface{}{
	Note{
		ID:        primitive.NewObjectID(),
		Title:     "First note",
		Body:      "Concurrency is not parallelism",
		CreatedAt: time.Now(),
	}, Note{
		ID:        primitive.NewObjectID(),
		Title:     "Second note",
		Body:      "A little copying is better than a little dependency",
		CreatedAt: time.Now(),
	}, Note{
		ID:        primitive.NewObjectID(),
		Title:     "Third note",
		Body:      "Don't communicate by sharing memory, share memory by communicating",
		CreatedAt: time.Now(),
	}, Note{
		ID:        primitive.NewObjectID(),
		Title:     "Fourth note",
		Body:      "Don't just check errors, handle them gracefully",
		CreatedAt: time.Now(),
	},
}

func main() {
	ctx := context.Background()
	clientOpts := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		fmt.Println(err)
		return
	}

	db := client.Database(DBName)
	coll := db.Collection(notesCollection)

	coll.Indexes().DropAll(ctx)

	// Options
	indexOptions := options.Index().SetUnique(true)
	indexKeys := bsonx.MDoc{
		"title": bsonx.Int32(-1),
	}

	noteIndexModel := mongo.IndexModel{
		Options: indexOptions,
		Keys:    indexKeys,
	}

	_, err = coll.Indexes().CreateOne(ctx, noteIndexModel)
	if err != nil {
		fmt.Println(err)
		return
	}

	textIndexModel := mongo.IndexModel{
		Options: options.Index().SetBackground(true),
		Keys: bsonx.MDoc{
			"title": bsonx.String("text"),
			"body":  bsonx.String("text"),
		},
	}

	_, err = coll.Indexes().CreateOne(ctx, textIndexModel)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = coll.InsertMany(ctx, notes)
	if err != nil {
		fmt.Println(err)
		return
	}

	n := Note{}
	fmt.Println("First example")
	cursor, err := coll.Find(ctx, bson.M{
		"$text": bson.M{
			"$search": "note",
		},
	})

	for cursor.Next(ctx) {
		cursor.Decode(&n)
		fmt.Println(n)
	}

	fmt.Println("Second example")
	cursor, err = coll.Find(ctx, bson.M{
		"$text": bson.M{
			"$search": "gracefully",
		},
	})

	for cursor.Next(ctx) {
		cursor.Decode(&n)
		fmt.Println(n)
	}
}
