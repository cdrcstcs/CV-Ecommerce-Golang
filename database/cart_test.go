package database
import (
	"context"
	"testing"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"ecommerce/models"
)
var (
	testProdCollection *mongo.Collection
	testUserCollection *mongo.Collection
)
func setupProductAndUser(t *testing.T, productID primitive.ObjectID, userID primitive.ObjectID) {
	product := models.ProductUser{
		Product_ID:    productID,
		Price: 100,
	}
	_, err := testProdCollection.InsertOne(context.Background(), product)
	require.NoError(t, err)
	user := models.User{
		ID:       userID,
		UserCart: []models.ProductUser{product},
	}
	_, err = testUserCollection.InsertOne(context.Background(), user)
	require.NoError(t, err)
}
func TestAddProductToCart(t *testing.T) {
	productID := primitive.NewObjectID()
	userID := primitive.NewObjectID().Hex()
	setupProductAndUser(t, productID, primitive.NewObjectID())
	err := AddProductToCart(context.Background(), testProdCollection, testUserCollection, productID, userID)
	require.NoError(t, err)
	var updatedUser models.User
	err = testUserCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&updatedUser)
	require.NoError(t, err)
	assert.Len(t, updatedUser.UserCart, 2) 
}
func TestRemoveCartItem(t *testing.T) {
	productID := primitive.NewObjectID()
	userID := primitive.NewObjectID().Hex()
	setupProductAndUser(t, productID, primitive.NewObjectID())
	err := RemoveCartItem(context.Background(), testProdCollection, testUserCollection, productID, userID)
	require.NoError(t, err)
	var updatedUser models.User
	err = testUserCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&updatedUser)
	require.NoError(t, err)
	assert.Empty(t, updatedUser.UserCart)
}
func TestBuyItemFromCart(t *testing.T) {
	productID := primitive.NewObjectID()
	userID := primitive.NewObjectID().Hex()
	setupProductAndUser(t, productID, primitive.NewObjectID())
	err := BuyItemFromCart(context.Background(), testUserCollection, userID)
	require.NoError(t, err)
	var updatedUser models.User
	err = testUserCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&updatedUser)
	require.NoError(t, err)
	assert.Empty(t, updatedUser.UserCart)
	assert.Len(t, updatedUser.Order_Status, 1)
}
func TestInstantBuyer(t *testing.T) {
	productID := primitive.NewObjectID()
	userID := primitive.NewObjectID().Hex()
	setupProductAndUser(t, productID, primitive.NewObjectID())

	err := InstantBuyer(context.Background(), testProdCollection, testUserCollection, productID, userID)
	require.NoError(t, err)

	var updatedUser models.User
	err = testUserCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&updatedUser)
	require.NoError(t, err)
	assert.Len(t, updatedUser.Order_Status, 1)
	assert.Equal(t, productID, updatedUser.Order_Status[0].Order_ID)
}