package controllers

import (
	"Nice-Things-Backend/initializers"
	"Nice-Things-Backend/models"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUp (c *gin.Context) {
	// Get email, first name, last name, and if they are an admin

	var body struct {
		Email string
		FirstName string
		LastName string
		Admin bool
		Sender int
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})

		return
	}

	fmt.Printf("SIGNUP BODY: %+v\n", body)

	var currentUser models.User
	initializers.DB.First(&currentUser, "id = ?", body.Sender)

	if currentUser.ID != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Not admin",
		})
	}

	// Create password = lower(LastName + FirstName)
	password := strings.ToLower(body.LastName + body.FirstName)

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to hash a password",
		})

		return
	}

	// Create User
	user := models.User{
		Email: body.Email, 
		FirstName: body.FirstName, 
		LastName: body.LastName, 
		Password: string(hash), 
		Admin: body.Admin,
	}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to create user",
		})

		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"Default Password": password,
	})
}

func SignIn (c *gin.Context) {
	// Get the email and pass off req body
	var body struct {
		Email string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})

		return
	}

	// Look up requested user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid email",
		})

		return
	}

	// Compare submitted password with saved user password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid password",
		})

		return
	}

	// Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to create token",
		})

		return
	}

	// send it back
	c.SetCookie("Authorization", tokenString, 3600 * 24 * 30, "" , "localhost:3000", true, true)
	c.Writer.Header().Set("Set-Cookie", "Authorization="+tokenString+"; SameSite=None; Secure")

	c.JSON(http.StatusOK, gin.H{
		"ID": user.ID,
		"Email": user.Email,
		"Token": tokenString,
		"Admin": user.Admin,
		"Name": user.FirstName + " " + user.LastName,
	})
}	


func ChangePassword (c *gin.Context) {

	var body struct {
		OldPassword string
		NewPassword string
		Sender int
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})

		return
	}

	fmt.Printf("ChangePassword BODY: %+v\n", body)

	// Look up requested user
	var currentUser models.User
	initializers.DB.First(&currentUser, "id = ?", body.Sender)


	err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(body.OldPassword))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid password",
		})

		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to hash a password",
		})

		return
	}

	initializers.DB.Model(&currentUser).Update("Password", string(hash))

	c.JSON(200, gin.H{})
}
