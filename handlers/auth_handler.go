package handlers

import (
	"blogpost/config"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type SignUpRequest struct {
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"firstname" validate:"required"`
	LastName  string `json:"lastname" validate:"required"`
	Role      string `json:"role" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func SignUp(c *fiber.Ctx) error {
	var req SignUpRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Insufficient Details, Please try Again!"})
	}

	userExists, err := userExistsByEmail(req.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Error checking user existence"})
	}
	if userExists {
		return c.Status(400).JSON(fiber.Map{"message": "Email already exists"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Error in hashing password"})
	}

	user := bson.M{
		"username":  req.Username,
		"email":     req.Email,
		"password":  string(hashedPassword),
		"firstname": req.FirstName,
		"lastname":  req.LastName,
		"role":      req.Role,
	}

	collection := config.DB.Collection("users")
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Error saving user to database"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Signup Successful"})
}

func userExistsByEmail(email string) (bool, error) {
	var result bson.M
	collection := config.DB.Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid request body"})
	}

	var user bson.M
	collection := config.DB.Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(401).JSON(fiber.Map{"message": "Invalid email or password"})
		}
		return c.Status(500).JSON(fiber.Map{"message": "Error finding user"})
	}

	if password, ok := user["password"].(string); ok {
		if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password)); err != nil {
			return c.Status(401).JSON(fiber.Map{"message": "Invalid credentials"})
		}
	} else {
		return c.Status(400).JSON(fiber.Map{"message": "Error retrieving password"})
	}
	token, err := config.CreateToken(user["_id"].(primitive.ObjectID).Hex())
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Error in Session token creation...!"})
	}
	return c.Status(200).JSON(fiber.Map{"message": "Login Successful", "token": token})
}
