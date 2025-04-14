package main

import (
    "time"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)

//defining a model for the table users
type User struct {
	ID uint   `gorm:"primaryKey;autoIncrement"`
	Username  string  `gorm:"unique;not null"`
	Password string  `gorm:"not null"`
}

var db *gorm.DB
var err error

func main() {

	//connecting to the database

	dsn := "  user=postgres password=051211308 dbname=spring_camp port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}

	fmt.Println("Successfully connected to the database!")
    
    //creating a table
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("Failed to auto migrate the schema")
	}

	

	

    //creating router
	router := gin.Default()
     

	//------------------------------------ROUTES------------------------------------------

	//create a user
	router.POST("/register", CreateUser)


	//login 
	router.POST("/login", LoginUser)


	//running server
	router.Run(":8080")
}


//------------------------------------HANDLER FUNCTIONS-----------------------------------

func CreateUser(c* gin.Context){
    var newUser User 

	if err := c.BindJSON(&newUser); err != nil {
        return
    }


    //hashing the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}
		newUser.Password = string(hashedPassword)

	result := db.Create(&newUser)
	if result.Error != nil {
		panic("Failed to create user")
	}
	
}
//TODO: See if a user with the same name already exists
func LoginUser(c *gin.Context){
	var inputUser User
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// Fetch the user from the database
	var dbUser User
	if err := db.Where("username = ?", inputUser.Username).First(&dbUser).Error; err != nil {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	// Compare passwords
	if err := comparePasswords(dbUser.Password, inputUser.Password); err != nil {
		c.JSON(401, gin.H{"error": "Invalid password"})
		return
	}

	// Create a JWT
	token, err := createToken(dbUser.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create token"})
		return
	}

	// Set the JWT in a cookie 
	c.SetCookie("token", token, 3600, "/", "localhost", false, true)

	c.JSON(200, gin.H{"message": "Login successful", "token":token})
}


//comparing hashed and entered password
func comparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

//setting jwt secret key
var jwtKey = []byte("whydidthechickencrosstheroad") 

func createToken(username string) (string, error) {
	// Define the token claims
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
