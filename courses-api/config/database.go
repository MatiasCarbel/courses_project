package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB() (*mongo.Database, error) {
    var client *mongo.Client
    var err error
    
    for i := 0; i < 5; i++ {
        clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        client, err = mongo.Connect(ctx, clientOptions)
        cancel()
        
        if err == nil {
            err = client.Ping(context.Background(), nil)
            if err == nil {
                break
            }
        }
        
        log.Printf("Failed to connect to MongoDB (attempt %d/5): %v", i+1, err)
        time.Sleep(2 * time.Second)
    }
    
    if err != nil {
        return nil, err
    }

    return client.Database("coursesdb"), nil
} 