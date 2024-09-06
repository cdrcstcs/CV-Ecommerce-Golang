package database
import (
	"context"
	"os"
	"testing"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
var testClient *mongo.Client
func TestMain(m *testing.M) {
	err := godotenv.Load("D:/CV-Projects/MainCV/CV-Ecommerce-Golang/.env")
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}
	mongoURI := os.Getenv("MONGO")
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic("Failed to connect to MongoDB: " + err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		panic("Failed to ping MongoDB: " + err.Error())
	}
	testClient = client
	code := m.Run()
	client.Disconnect(context.Background())
	os.Exit(code)
}
func TestDBSet(t *testing.T) {
	client := DBSet()
	require.NotNil(t, client, "DBSet should return a non-nil client")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := client.Ping(ctx, nil)
	require.NoError(t, err, "Failed to ping MongoDB")
}
func TestUserData(t *testing.T) {
	collection := UserData(testClient, "users")
	require.NotNil(t, collection, "UserData should return a non-nil collection")
	count, err := collection.CountDocuments(context.Background(), bson.M{})
	require.NoError(t, err, "Error counting documents in the collection")
	assert.True(t, count >= 0, "Collection should be accessible")
}
func TestProductData(t *testing.T) {
	collection := ProductData(testClient, "products")
	require.NotNil(t, collection, "ProductData should return a non-nil collection")
	count, err := collection.CountDocuments(context.Background(), bson.M{})
	require.NoError(t, err, "Error counting documents in the collection")
	assert.True(t, count >= 0, "Collection should be accessible")
}