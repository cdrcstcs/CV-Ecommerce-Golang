package token

import (
	"context"
	"log"
	"os"
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var (
	client         *mongo.Client
	mockDB         *mongo.Database
	mockUserColl   *mongo.Collection
)
func setup() {
	ctx := context.Background()
	err := godotenv.Load("D:/CV-Projects/MainCV/CV-Ecommerce-Golang/.env") 
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mongoUri := os.Getenv("MONGO")
	SECRET_KEY = os.Getenv("SECRET_LOVE") 
	clientOptions := options.Client().ApplyURI(mongoUri)
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	mockDB = client.Database("testdb")
	mockUserColl = mockDB.Collection("Users")
	UserData = mockUserColl
}
func teardown() {
	ctx := context.Background()
	if err := mockDB.Drop(ctx); err != nil {
		log.Fatal(err)
	}
	if err := client.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
}
func TestTokenGeneration(t *testing.T) {
	setup()
	defer teardown()
	email := "test@example.com"
	firstname := "John"
	lastname := "Doe"
	uid := "123456"
	token, refreshtoken, err := TokenGenerator(email, firstname, lastname, uid)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, refreshtoken)
}
func TestTokenValidation(t *testing.T) {
	setup()
	defer teardown()
	email := "test@example.com"
	firstname := "John"
	lastname := "Doe"
	uid := "123456"
	token, _, err := TokenGenerator(email, firstname, lastname, uid)
	assert.NoError(t, err)
	claims, msg := ValidateToken(token)
	assert.Empty(t, msg)
	assert.NotNil(t, claims)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, firstname, claims.First_Name)
	assert.Equal(t, lastname, claims.Last_Name)
	assert.Equal(t, uid, claims.Uid)
	invalidToken := token + "invalid"
	_, msg = ValidateToken(invalidToken)
	assert.Equal(t, "signature is invalid", msg)
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(-time.Hour).Unix(),
		},
	})
	tokenString, _ := expiredToken.SignedString([]byte(SECRET_KEY))
	_, msg = ValidateToken(tokenString)
	assert.Contains(t, msg, "token is expired")
}
func TestUpdateAllTokens(t *testing.T) {
	setup()
	defer teardown()
	userID := primitive.NewObjectID().Hex()
	_, err := mockUserColl.InsertOne(context.Background(), bson.M{
		"user_id": userID,
	})
	if err != nil {
		log.Fatal(err)
	}
	token, refreshtoken, err := TokenGenerator("test@example.com", "John", "Doe", userID)
	assert.NoError(t, err)
	UpdateAllTokens(token, refreshtoken, userID)
	var result bson.M
	err = mockUserColl.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, token, result["token"])
	assert.Equal(t, refreshtoken, result["refresh_token"])
}