package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var tableName = "Users"
var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var db = dynamodb.New(sess)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type VisibleUser struct {
	Username string `json:"username"`
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatalf("Error when hashing pass %s", err)
	}
	return string(bytes)
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func main() {
	router := gin.Default()
	router.GET("/users", GetUsers)
	router.GET("/users/:id", GetUserById)
	router.POST("/users", AddUser)
	router.DELETE("/users/:id", DeleteUserById)
	router.Run("localhost:8080")
}

func DeleteUserById(c *gin.Context) {
	userID := c.Param("id")
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(userID),
			},
		},
	}

	_, err := db.DeleteItem(input)
	if err != nil {
		log.Fatalf("Error when deleting item %s", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func GetUsers(c *gin.Context) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName), // replace with your table name
	}

	result, err := db.Scan(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving items from DB"})
		return
	}

	users := make([]VisibleUser, len(result.Items))
	for i, item := range result.Items {
		user := VisibleUser{}
		err = dynamodbattribute.UnmarshalMap(item, &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error unmarshalling item"})
			return
		}
		users[i] = user
	}

	c.JSON(http.StatusOK, users)
}

func GetUserById(c *gin.Context) {
	id := c.Param("id")
	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		log.Fatalf("Error when calling GetItem %s", err)
	}

	if result.Item == nil {
		log.Fatalf("Couldn't find user with that ID")
	}

	user := VisibleUser{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		log.Fatalf("Failed to unmarshal record")
	}

	c.JSON(http.StatusOK, user)
}

func AddUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := uuid.New().String()
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: &userID,
			},
			"Username": {
				S: aws.String(user.Username),
			},
			"Password": {
				S: aws.String(HashPassword(user.Password)),
			},
		},
	}

	_, err := db.PutItem(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting item into DB"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
