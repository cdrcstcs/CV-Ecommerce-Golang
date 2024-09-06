package token
import (
	"context"
	"log"
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/docker/go-connections/nat"
)
var (
	client         *mongo.Client
	mockDB         *mongo.Database
	mockUserColl   *mongo.Collection
	mongoContainer testcontainers.Container
	mongoPort      string
)
func setup() {
	ctx := context.Background()
	mongoContainerPort := "27017/tcp"
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{mongoContainerPort},
		WaitingFor:   wait.ForListeningPort(nat.Port(mongoContainerPort)),
	}
	var err error
	mongoContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:           true,
	})
	if err != nil {
		log.Fatal(err)
	}
	mappedPort, err := mongoContainer.MappedPort(ctx, nat.Port(mongoContainerPort))
	if err != nil {
		log.Fatal(err)
	}
	mongoPort = mappedPort.Port()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:" + mongoPort)
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	mockDB = client.Database("testdb")
	mockUserColl = mockDB.Collection("Users")
	UserData = mockUserColl
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	SECRET_KEY = "test_secret_key"
}
func teardown() {
	ctx := context.Background()
	if err := client.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
	if err := mongoContainer.Terminate(ctx); err != nil {
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
	assert.Equal(t, "token is expired", msg)
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