package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Title  string             `bson:"title,omitempty"`
	Author string             `bson:"author,omitempty"`
	Year   int                `bson:"year_published,omitempty"`
}

func createBookRepository(uri, db, col string) *repository {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return &repository{collection: client.Database(db).Collection(col)}
}

type repository struct {
	collection *mongo.Collection
}

func (r *repository) createBook(book Book) (*mongo.InsertOneResult, error) {
	return r.collection.InsertOne(context.TODO(), book)
}

func (r *repository) readBook(id interface{}) *mongo.SingleResult {
	filter := bson.M{"_id": id}
	return r.collection.FindOne(context.TODO(), filter)
}

func (r *repository) updateBook(id interface{}, book Book) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": book}
	return r.collection.UpdateOne(context.TODO(), filter, update)
}

func (r *repository) deleteBook(id interface{}) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": id}
	return r.collection.DeleteMany(context.TODO(), filter)
}

func main() {
	uri := "mongodb+srv://admin:admin@cluster0.xtwwu.mongodb.net"
	database := "myDB"
	collection := "favorite_books"
	repo := createBookRepository(uri, database, collection)

	result, err := repo.createBook(Book{
		Title:  "Invisible Cities",
		Author: "Italo Calvino",
		Year:   1974,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)

	book := repo.readBook(result.InsertedID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", book)

	updateResult, err := repo.updateBook(result.InsertedID, Book{
		Title:  "Bumi manusia",
		Author: "Pramoedya Ananta Toer",
		Year:   1980,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Documents matched: %v\n", updateResult.MatchedCount)
	fmt.Printf("Documents updated: %v\n", updateResult.ModifiedCount)

	book = repo.readBook(result.InsertedID)
	fmt.Printf("%v\n", book)

	deleteResult, err := repo.deleteBook(result.InsertedID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Number of documents deleted: %d\n", deleteResult.DeletedCount)
}
