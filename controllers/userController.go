package controllers

import (
	"strconv"
	"time"

	"github.com/alpha_batta/database"
	"github.com/alpha_batta/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"

// Register function
func Register(r *fiber.Ctx) error {

	/*
		Creating a data variable which is a map to hold the user values
	*/
	var data map[string]string

	/*
		Parsing the data to check for errors in the inputs
	*/
	err := r.BodyParser(&data)

	if err != nil {
		return err
	}

	// Hashing the password
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		UserName: data["username"],
		Email:    data["email"],
		Password: string(password),
	}

	database.Database.Db.Create(&user)

	return r.JSON(user)
}

// Login function
func Login(r *fiber.Ctx) error {
	var data map[string]string
	if err := r.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	/*
		checking for the user in the database by email
	*/
	database.Database.Db.Where("email = ?", data["email"]).First(&user)

	if user.ID == 0 {
		r.Status(fiber.StatusNotFound)
		return r.JSON(fiber.Map{"message": "User not found"})
	}

	/*
		Comparing the entered password and the password stored in the database
	*/
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		r.Status(fiber.StatusBadRequest)
		return r.JSON(fiber.Map{"message": "Please enter a valid password"})
	}

	/*
		Generating standards claims then using them to generate a token
	*/

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		r.Status(fiber.StatusInternalServerError)
		return r.JSON(fiber.Map{"message": "could not login"})
	}

	/*
		Storing the token in a cookie
	*/
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	r.Cookie(&cookie)

	return r.JSON(fiber.Map{"message": "Success"})
}

// User function to retreive logged in user

func User(r *fiber.Ctx) error {
	// getting the cookie by name
	cookie := r.Cookies("jwt")

	// retreiving the token
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		r.Status(fiber.StatusUnauthorized)
		return r.JSON(fiber.Map{"message": "unauthenticated"})
	}

	// Getting the claims from the token
	claims := token.Claims.(*jwt.StandardClaims)

	// retreiving the user by claims and returning it

	var user models.User

	database.Database.Db.Where("id = ?", claims.Issuer).First(&user)
	return r.JSON(user)
}

// Log out function

func Logout(r *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	r.Cookie(&cookie)

	return r.JSON(fiber.Map{"message": "Success"})
}
