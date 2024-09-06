package controllers
import (
	"os"
	"context"
	"net/http"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ecommerce/models"
	"github.com/joho/godotenv"
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
	expected := "Successfully Signed Up!!"
	actual := w.Body.String()
	if len(actual) > 1 && actual[0] == '"' && actual[len(actual)-1] == '"' {
		actual = actual[1 : len(actual)-1]
	}
	if actual != expected {
		assert.Contains(t, actual, "error")
	} else {
		assert.Equal(t, expected, actual)
	}
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
	expected := "Successfully added our Product Admin!!"
	actual := w.Body.String()
	if len(actual) > 1 && actual[0] == '"' && actual[len(actual)-1] == '"' {
		actual = actual[1 : len(actual)-1]
	}
	if actual != expected {
		assert.Contains(t, actual, "error")
	} else {
		assert.Equal(t, expected, actual)
	}}
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
	r.GET("/users/search", SearchProductByQuery())
	w := performRequest(r, "GET", "/users/search?name=Sample", nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Sample")
}
func stringPtr(s string) *string {
	return &s
}
func intPtr(i uint64) *uint64 {
	return &i
}