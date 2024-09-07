package database
import (
	"context"
	"os"
	"testing"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ecommerce/models"
)
var (
	client         *mongo.Client
	mockDB         *mongo.Database
	mockUserColl   *mongo.Collection
	mockProdColl   *mongo.Collection
)
func setup() {
	ctx := context.Background()
	err := godotenv.Load("D:/CV-Projects/MainCV/CV-Ecommerce-Golang/.env")
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}
	mongoURI := os.Getenv("MONGO")
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	mockDB = client.Database("testdb")
	mockUserColl = mockDB.Collection("users")
	mockProdColl = mockDB.Collection("products")
	mockUserColl.Drop(ctx)
	mockProdColl.Drop(ctx)
}
func teardown() {
	ctx := context.Background()
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}
func setupProductAndUser(t *testing.T, productID primitive.ObjectID, userID primitive.ObjectID) {
	product := models.ProductUser{
		Product_ID: productID,
		Price:      100,
	}
	_, err := mockProdColl.InsertOne(context.Background(), product)
	require.NoError(t, err)
	user := models.User{
		ID:       userID,
		UserCart: []models.ProductUser{product},
	}
	_, err = mockUserColl.InsertOne(context.Background(), user)
	require.NoError(t, err)
}
func TestAddProductToCart(t *testing.T) {
	setup()
	defer teardown()
	productID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	setupProductAndUser(t, productID, userID)
	err := AddProductToCart(context.Background(), mockProdColl, mockUserColl, productID, userID.Hex())
	require.NoError(t, err)
	var updatedUser models.User
	err = mockUserColl.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&updatedUser)
	require.NoError(t, err)
	assert.Len(t, updatedUser.UserCart, 2)
}
func TestRemoveCartItem(t *testing.T) {
	setup()
	defer teardown()
	productID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	setupProductAndUser(t, productID, userID)
	err := RemoveCartItem(context.Background(), mockProdColl, mockUserColl, productID, userID.Hex())
	require.NoError(t, err)
	var updatedUser models.User
	err = mockUserColl.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&updatedUser)
	require.NoError(t, err)
	assert.Empty(t, updatedUser.UserCart)
}
func TestBuyItemFromCart(t *testing.T) {
	setup()
	defer teardown()
	productID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	setupProductAndUser(t, productID, userID)
	err := BuyItemFromCart(context.Background(), mockUserColl, userID.Hex())
	require.NoError(t, err)
	var updatedUser models.User
	err = mockUserColl.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&updatedUser)
	require.NoError(t, err)
	assert.Empty(t, updatedUser.UserCart)
}
func TestInstantBuyer(t *testing.T) {
	setup()
	defer teardown()
	productID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	setupProductAndUser(t, productID, userID)
	err := InstantBuyer(context.Background(), mockProdColl, mockUserColl, productID, userID.Hex())
	require.NoError(t, err)
	var updatedUser models.User
	err = mockUserColl.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&updatedUser)
	require.NoError(t, err)
	assert.Len(t, updatedUser.Order_Status, 0)
}