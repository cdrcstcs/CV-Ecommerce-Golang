package controllers
import (
	"context"
	"net/http"
	"time"
	"ecommerce/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
func AddAddress() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.Query("id")
        if userID == "" {
            c.Header("Content-Type", "application/json")
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
            return
        }
        addressID, err := primitive.ObjectIDFromHex(userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
            return
        }
        var address models.Address
        address.Address_id = primitive.NewObjectID()
        if err = c.BindJSON(&address); err != nil {
            c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
            return
        }
        ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()
        matchFilter := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: addressID}}}}
        unwind := bson.D{{Key: "$unwind", Value: "$address"}}
        group := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$_id"}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}
        cursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{matchFilter, unwind, group})
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
            return
        }
        defer cursor.Close(ctx) 
        var addressInfo []bson.M
        if err = cursor.All(ctx, &addressInfo); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing addresses"})
            return
        }
        var size int32
        if len(addressInfo) > 0 {
            count, ok := addressInfo[0]["count"].(int32)
            if !ok {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading address count"})
                return
            }
            size = count
        }
        if size < 2 {
            filter := bson.D{{Key: "_id", Value: addressID}}
            update := bson.D{{Key: "$push", Value: bson.D{{Key: "address", Value: address}}}}
            _, err := UserCollection.UpdateOne(ctx, filter, update)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating address"})
                return
            }
            c.JSON(http.StatusOK, gin.H{"message": "Address added successfully"})
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Address limit reached"})
        }
    }
}
func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid"})
			c.Abort()
			return
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, err)
		}
		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editaddress.House}, {Key: "address.0.street_name", Value: editaddress.Street}, {Key: "address.0.city_name", Value: editaddress.City}, {Key: "address.0.pin_code", Value: editaddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "Something Went Wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully Updated the Home address")
	}
}
func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Wrong id not provided"})
			c.Abort()
			return
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, err)
		}
		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editaddress.House}, {Key: "address.1.street_name", Value: editaddress.Street}, {Key: "address.1.city_name", Value: editaddress.City}, {Key: "address.1.pin_code", Value: editaddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "something Went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully updated the Work Address")
	}
}
func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		addresses := make([]models.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(404, "Wromg")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully Deleted!")
	}
}