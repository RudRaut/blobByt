package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// FileMetadata holds metadata info about uploaded files.
type FileMetadata struct {
	ID            string    `bson:"_id,omitempty"`
	BlobID        string    `bson:"blobID"`
	Name          string    `bson:"name"`
	Size          int64     `bson:"size"`
	FileType      string    `bson:"fileType"`
	EncryptionKey string    `bson:"encryptionKey"`
	Epochs        int       `bson:"epochs"`
	Description   string    `bson:"description,omitempty"`
	UploadTime    time.Time `bson:"uploadTime"`
}

// Client is a package-level MongoDB client accessible to all DB functions.
var Client *mongo.Client

// Connect initializes the MongoDB client, connects, pings, and saves it to package variable.
func Connect(uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOpts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return err
	}

	// Ping to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	fmt.Println("Successfully connected and pinged MongoDB!")

	Client = client
	return nil
}

// InsertFileMetadata inserts file metadata into the "files" collection.
// Uses package-level Client, no need to pass client explicitly.
func InsertFileMetadata(metadata FileMetadata) error {
	if Client == nil {
		return fmt.Errorf("mongo client not initialized")
	}

	collection := Client.Database("walrus").Collection("files")

	// Set upload time here
	metadata.UploadTime = time.Now()

	metadata.ID = metadata.BlobID

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, metadata)
	return err
}

// FindByBlobID retrieves a FileMetadata document by the blobID (used as _id).
func GetMetadataByBlobID(blobID string) (*FileMetadata, error) {
	if Client == nil {
		return nil, fmt.Errorf("mongo client not initialized")
	}

	collection := Client.Database("walrus").Collection("files")

	var metadata FileMetadata
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"_id": blobID}).Decode(&metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}
