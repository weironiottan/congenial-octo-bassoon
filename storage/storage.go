// Package storage contains the code to persist and retrieve orders from a
// database
package storage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"

	"github.com/levenlabs/go-llog"
)

// Instance holds a database connection for use in the storage methods
type Instance struct {
	// database is the name of the database within the storage engine and being a
	// variable we'll randomize this in tests so we don't need to wipe the database
	// between every test run
	// you should use this as the database name for all of your methods to simplify
	// testing
	database   string
	db         *mongo.Client
	collection *mongo.Collection

	// this is where you'd store any database connections like a *mongo.Client or
	// *sql.DB

}

func New(overrideDatabase string) *Instance {
	// create a pointer to an Instance that we will return after initialization
	inst := &Instance{}
	// if they sent overrideDatabase then use that, like for tests, otherwise
	// fallback to a static name for production, staging, etc
	// you might have some environmental variables or flags instead to do this but
	// for this project this is sufficient
	if overrideDatabase != "" {
		inst.database = overrideDatabase
	} else {
		inst.database = "order_up"
	}

	// TODO: code for connecting to the database and storing the connected driver
	db, err := inst.openDB()
	if err != nil {
		log.Fatal(err)
	}
	inst.db = db
	inst.collection = db.Database("test").Collection("Order")
	// normally I would include a disconnect somewhere in the main, since we do not
	// want to open/close the connection for each connection. I will leave it out because IDK where to put it here

	// instance on inst

	// give the ensureSchema function only 15 seconds to complete
	// after 15 seconds the context will return DeadlineExceeded errors which should
	// cause any functions downstream to error out
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// if we don't call cancel then the ctx will leak so we make sure that cancel
	// is called no matter what when we're done
	defer cancel()
	// we want to make sure the database is ready to accept requests and if that
	// fails we need to fatal
	if err := inst.ensureSchema(ctx); err != nil {
		llog.Fatal("failed to ensure schema", llog.ErrKV(err))
	}
	return inst
}

func (i *Instance) ensureSchema(ctx context.Context) error {
	// TODO: this is where you'll do any schema setup (CREATE DATABASE or CREATE
	// TABLE), if necessary, and since this will be called every time the service
	// starts or every time you run tests, it should not fail if the schema is
	// already setup
	// for example you might need to add a unique index on the order's ID field ;)
	return nil
}

func (i *Instance) openDB() (*mongo.Client, error) {

	// Set connection options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	client.Database("edmEvents").Collection("lasVegasEdmEventsCollection")

	fmt.Println("Connected to MongoDB!")
	return client, nil
}
