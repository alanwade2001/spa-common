package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"k8s.io/klog/v2"
)

// MongoService s
type MongoService struct {
	uriTemplate string
	username    string
	password    string
	database    string
	collection  string
	timeout     time.Duration
	registry    *bsoncodec.Registry
}

// NewMongoService s
func NewMongoService(uriTemplate string, user string, password string, database string, collection string, timeout time.Duration, registry *bsoncodec.Registry) *MongoService {
	return &MongoService{
		uriTemplate: uriTemplate,
		username:    user,
		password:    password,
		database:    database,
		collection:  collection,
		timeout:     timeout,
		registry:    registry,
	}
}

// MongoConnection s
type MongoConnection struct {
	Client *mongo.Client
	Ctx    context.Context
	Cancel context.CancelFunc
}

// Disconnect f
func (mc *MongoConnection) Disconnect() {
	mc.Cancel()
	mc.Client.Disconnect(mc.Ctx)
}

// Connect f
func (ms MongoService) Connect() MongoConnection {

	connectionURI := fmt.Sprintf(ms.uriTemplate, ms.username, ms.password)

	// structcodec, _ := bsoncodec.NewStructCodec(bsoncodec.JSONFallbackStructTagParser)
	//reg := bson.NewRegistryBuilder().
	//	RegisterTypeEncoder(reflect.TypeOf(generated.CustomerModel{}), structcodec).
	//	RegisterTypeDecoder(reflect.TypeOf(generated.CustomerModel{}), structcodec).Build()

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI).SetRegistry(ms.registry))
	if err != nil {
		klog.Warningf("Failed to create client: %v", err)
	}

	connectTimeout := ms.timeout * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)

	err = client.Connect(ctx)
	if err != nil {
		klog.Warningf("Failed to connect to cluster: %v", err)
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		klog.Warningf("Failed to ping cluster: %v", err)
	}

	klog.Infof("Connected to MongoDB!")

	return MongoConnection{client, ctx, cancel}
}

func (ms MongoService) GetCollection(connection MongoConnection) *mongo.Collection {
	return connection.Client.Database(ms.database).Collection(ms.collection)
}
