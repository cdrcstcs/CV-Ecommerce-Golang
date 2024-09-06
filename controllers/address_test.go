package controllers
import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"ecommerce/models"
)
func TestAddAddress(t *testing.T) {
	setup()
	defer teardown()
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/address/add", AddAddress())
	userID := primitive.NewObjectID().Hex()
	_, _ = mockUserColl.InsertOne(context.Background(), models.User{
		ID:        primitive.NewObjectID(),
		User_ID:   userID,
		Address_Details:   []models.Address{},
	})
	address := models.Address{
		House:     stringPtr("123"),
		Street:    stringPtr("Main St"),
		City:      stringPtr("Cityville"),
		Pincode:   stringPtr("12345"),
	}
	w := performRequest(r, "POST", "/address/add?id="+userID, address)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Successfully Added Address", w.Body.String())
}
func TestEditHomeAddress(t *testing.T) {
	setup()
	defer teardown()
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PUT("/address/home", EditHomeAddress())
	userID := primitive.NewObjectID().Hex()
	_, _ = mockUserColl.InsertOne(context.Background(), models.User{
		ID:        primitive.NewObjectID(),
		User_ID:   userID,
		Address_Details:   []models.Address{
			{
				House: 	stringPtr("123"), 
				Street: stringPtr("Old St"), 
				City: 	stringPtr("Oldville"), 
				Pincode:stringPtr("54321"),
			}},
	})
	updatedAddress := models.Address{
		House:   stringPtr("456"),
		Street:  stringPtr("New St"),
		City:    stringPtr("Newville"),
		Pincode: stringPtr("67890"),
	}
	w := performRequest(r, "PUT", "/address/home?id="+userID, updatedAddress)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Successfully Updated the Home address", w.Body.String())
}
func TestEditWorkAddress(t *testing.T) {
	setup()
	defer teardown()
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PUT("/address/work", EditWorkAddress())
	userID := primitive.NewObjectID().Hex()
	_, _ = mockUserColl.InsertOne(context.Background(), models.User{
		ID:        primitive.NewObjectID(),
		User_ID:   userID,
		Address_Details:   []models.Address{
			{
				House: 	stringPtr("789"), 
				Street: stringPtr("Office St"), 
				City: 	stringPtr("Worktown"), 
				Pincode:stringPtr("98765"),
			}, 
			{
				House: 	stringPtr("321"), 
				Street: stringPtr("Old Office St"), 
				City: 	stringPtr("Worktown"), 
				Pincode:stringPtr("56789"),
			}},
	})
	updatedAddress := models.Address{
		House:   stringPtr("654"),
		Street:  stringPtr("New Office St"),
		City:    stringPtr("Newtown"),
		Pincode: stringPtr("12345"),
	}
	w := performRequest(r, "PUT", "/address/work?id="+userID, updatedAddress)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Successfully updated the Work Address", w.Body.String())
}
func TestDeleteAddress(t *testing.T) {
	setup()
	defer teardown()
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.DELETE("/address/delete", DeleteAddress())
	userID := primitive.NewObjectID().Hex()
	_, _ = mockUserColl.InsertOne(context.Background(), models.User{
		ID:        primitive.NewObjectID(),
		User_ID:   userID,
		Address_Details:   []models.Address{
			{
				House: 	stringPtr("789"), 
				Street: stringPtr("Office St"), 
				City: 	stringPtr("Worktown"), 
				Pincode:stringPtr("98765"),
			}},
	})
	w := performRequest(r, "DELETE", "/address/delete?id="+userID, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Successfully Deleted!", w.Body.String())
}
func performRequest(r *gin.Engine, method, url string, body interface{}) *httptest.ResponseRecorder {
	var requestBody *bytes.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		requestBody = bytes.NewReader(b)
	} else {
		requestBody = bytes.NewReader([]byte{})
	}
	req := httptest.NewRequest(method, url, requestBody)
	if method == "GET" {
		req = httptest.NewRequest(method, url, nil)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}