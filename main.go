package main

import (
	"blogpost/config"
	"blogpost/routes"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func init() {
	fmt.Println("Init called")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error in Loading .env File")
	}
	log.Println(".env is loaded Successfully")
	err = config.ConnectDB()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("DB Connected Succeessfully")
}

func main() {

	app := fiber.New()
	PORT := os.Getenv("PORT")
	address := fmt.Sprintf(":%s", PORT)
	// allowedOrigins:=[]string {"http://localhost:3000","http://localhost:5000"}
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))

	routes.SetupBlogRoutes(app)
	routes.SetupAuthRouter(app)
	err := app.Listen(address)
	if err != nil {
		log.Println("Error in starting server", err)
	}
	log.Println("Server Started on PORT: ", PORT)

}
