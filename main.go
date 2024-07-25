package main

import (
	"log"
	"os"
	"fmt"
	"github.com/akhil/ecommerce-yt/controllers"
	"github.com/akhil/ecommerce-yt/database"
	"github.com/akhil/ecommerce-yt/middleware"
	"github.com/akhil/ecommerce-yt/routes"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load()
    if err != nil {
        fmt.Println("Error loading .env file")
    }
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/listcart", controllers.GetItemFromCart())
	router.POST("/addaddress", controllers.AddAddress())
	router.PUT("/edithomeaddress", controllers.EditHomeAddress())
	router.PUT("/editworkaddress", controllers.EditWorkAddress())
	router.GET("/deleteaddresses", controllers.DeleteAddress())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())
	log.Fatal(router.Run(":" + port))
}
