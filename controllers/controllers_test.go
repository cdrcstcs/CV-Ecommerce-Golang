package controllers
import (
	"context"
	"net/http"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"ecommerce/models"
	"github.com/docker/go-connections/nat"
)
var (
	client          *mongo.Client
	mockDB          *mongo.Database
	mockUserColl    *mongo.Collection
	mockProdColl    *mongo.Collection
	mongoContainer  testcontainers.Container
	mongoPort       string
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
		panic(err)
	}
	mappedPort, err := mongoContainer.MappedPort(ctx, nat.Port(mongoContainerPort))
	if err != nil {
		panic(err)
	}
	mongoPort = mappedPort.Port()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:" + mongoPort)
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	mockDB = client.Database("testdb")
	mockUserColl = mockDB.Collection("users")
	mockProdColl = mockDB.Collection("products")
	UserCollection = mockUserColl
	ProductCollection = mockProdColl
}
func teardown() {
	ctx := context.Background()
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
	if err := mongoContainer.Terminate(ctx); err != nil {
		panic(err)
	}
}
func TestSignUp(t *testing.T) {
	setup()
	defer teardown()
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/signup", SignUp())
	user := models.User{
		Email:    stringPtr("test@example.com"),
		Password: stringPtr("password"),
		First_Name: stringPtr("John"),
		Last_Name:  stringPtr("Doe"),
		Phone:     stringPtr("1234567890"),
	}
	w := performRequest(r, "POST", "/signup", user)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "Successfully Signed Up!!", w.Body.String())
}
func TestLogin(t *testing.T) {
	setup()
	defer teardown()
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/login", Login())
	password := HashPassword("password")
	user := models.User{
		Email:    stringPtr("test@example.com"),
		Password: &password,
	}
	_, _ = mockUserColl.InsertOne(context.Background(), user)
	loginUser := models.User{
		Email:    stringPtr("test@example.com"),
		Password: stringPtr("password"),
	}
	w := performRequest(r, "POST", "/login", loginUser)
	assert.Equal(t, http.StatusFound, w.Code)
}
func TestProductViewerAdmin(t *testing.T) {
	setup()
	defer teardown()
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/product/admin", ProductViewerAdmin())
	product := models.Product{
		Product_Name: stringPtr("Sample Product"),
		Price:        intPtr(100),
	}
	w := performRequest(r, "POST", "/product/admin", product)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Successfully added our Product Admin!!", w.Body.String())
}
func TestSearchProduct(t *testing.T) {
	setup()
	defer teardown()
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/products", SearchProduct())
	product := models.Product{
		Product_Name: stringPtr("Sample Product"),
		Price:        intPtr(100),
	}
	_, _ = mockProdColl.InsertOne(context.Background(), product)
	w := performRequest(r, "GET", "/products", nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Sample Product")
}
func TestSearchProductByQuery(t *testing.T) {
	setup()
	defer teardown()
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/products/search", SearchProductByQuery())
	product := models.Product{
		Product_Name: stringPtr("Test Product"),
		Price:        intPtr(100),
	}
	_, _ = mockProdColl.InsertOne(context.Background(), product)
	w := performRequest(r, "GET", "/products/search?name=Test", nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Product")
}
func stringPtr(s string) *string {
	return &s
}
func intPtr(ui uint64) *uint64 {
	return &ui
}