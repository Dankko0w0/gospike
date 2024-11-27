package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	config     *Config
	client     *mongo.Client
	db         *mongo.Database
	connected  bool
	maxRetries int
}

func NewMongoDB(config *Config) *MongoDB {
	return &MongoDB{
		config:     config,
		maxRetries: 3,
	}
}

func (m *MongoDB) Connect(ctx context.Context) error {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		m.config.Username,
		m.config.Password,
		m.config.Host,
		m.config.Port,
	)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	m.client = client
	m.db = client.Database(m.config.Database)
	m.connected = true
	return nil
}

func (m *MongoDB) Disconnect(ctx context.Context) error {
	if m.client != nil {
		err := m.client.Disconnect(ctx)
		if err != nil {
			return fmt.Errorf("failed to disconnect from MongoDB: %v", err)
		}
		m.connected = false
	}
	return nil
}

func (m *MongoDB) Ping(ctx context.Context) error {
	if m.client == nil {
		return fmt.Errorf("database not connected")
	}
	return m.client.Ping(ctx, readpref.Primary())
}

func (m *MongoDB) IsConnected() bool {
	return m.connected
}

func (m *MongoDB) Reconnect(ctx context.Context) error {
	m.Disconnect(ctx)

	for i := 0; i < m.maxRetries; i++ {
		err := m.Connect(ctx)
		if err == nil {
			return nil
		}
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return fmt.Errorf("failed to reconnect after %d attempts", m.maxRetries)
}

// CRUD operations
func (m *MongoDB) Create(ctx context.Context, collection string, data interface{}) error {
	coll := m.db.Collection(collection)
	_, err := coll.InsertOne(ctx, data)
	return err
}

func (m *MongoDB) Read(ctx context.Context, collection string, filter interface{}, result interface{}) error {
	coll := m.db.Collection(collection)
	return coll.FindOne(ctx, filter).Decode(result)
}

func (m *MongoDB) Update(ctx context.Context, collection string, filter interface{}, update interface{}) error {
	coll := m.db.Collection(collection)
	_, err := coll.UpdateOne(ctx, filter, update)
	return err
}

func (m *MongoDB) Delete(ctx context.Context, collection string, filter interface{}) error {
	coll := m.db.Collection(collection)
	_, err := coll.DeleteOne(ctx, filter)
	return err
}

func (m *MongoDB) List(ctx context.Context, collection string, filter interface{}, results interface{}) error {
	coll := m.db.Collection(collection)
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return err
	}
	return cursor.All(ctx, results)
}
