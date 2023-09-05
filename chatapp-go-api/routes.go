package main

import (
	"example/chatapp-go-api/pkg/websocket"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func wshandler(w http.ResponseWriter, r *http.Request, pool *websocket.Pool, username string) {
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		log.Fatalf("Error when upgrading HTTP connection to Websocket protocol %s", err)
	}

	uniqueID := uuid.New().String()
	client := &websocket.Client{
		Username: username,
		Conn:     conn,
		Pool:     pool,
		Send:     make(chan websocket.SocketEvent),
		UserID:   uniqueID,
	}

	pool.Register <- client
	go client.Write()
	go client.Read()

}

// HTTP REST API for user authentication stuff w/ DynamoDB
var tableName = "Users"
var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var db = dynamodb.New(sess)

func Login(c *gin.Context) {
	var loginData User
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Username": {
				S: aws.String(loginData.Username),
			},
		},
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Error fetching user from DB"})
		return
	}

	if result.Item == nil {
		c.JSON(400, gin.H{"error": "Invalid username"})
		return
	}

	var dbUser User
	err = dynamodbattribute.UnmarshalMap(result.Item, &dbUser)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error unmarshalling DB result"})
		return
	}

	if !CheckPasswordHash(loginData.Password, dbUser.Password) {
		c.JSON(401, gin.H{"error": "Invalid password"})
		return
	}

	c.JSON(200, gin.H{"message": "Logged in successfully", "username": dbUser.Username})

}

func HomeHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello world")
}

func DeleteUserByUsername(c *gin.Context) {
	username := c.Param("username")
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Username": {
				S: aws.String(username),
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
		TableName: aws.String(tableName),
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

func AddUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]*dynamodb.AttributeValue{
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
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
