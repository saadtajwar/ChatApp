package main

import (
	"fmt"
	"log"
	"net/http"

	"example/chatapp-go-api/pkg/websocket"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type VisibleUser struct {
	Username string `json:"username"`
}

var tableName = "Users"
var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var db = dynamodb.New(sess)

func wshandler(w http.ResponseWriter, r *http.Request, pool *websocket.Pool) {
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		log.Fatalf("Error when upgrading HTTP connection to Websocket protocol %s", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()

}

func main() {
	router := gin.Default()
	pool := websocket.NewPool()
	go pool.Start()
	router.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request, pool)
	})
	router.GET("/", HomeHandler)
	router.GET("/users", GetUsers)
	router.POST("/users", AddUser)
	router.DELETE("/users/:username", DeleteUserByUsername)
	router.POST("/login", Login)
	router.Run(":8080")
}

func Login(c *gin.Context) {
	var loginData User
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
