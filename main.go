package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	service "meetcode-backend/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

type Message struct {
	Text string `json:"text"`
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan Message)
)

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Add client to map
	clients[conn] = true

	// Infinite loop to read incoming messages
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			delete(clients, conn)
			return
		}

		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Get next message from broadcast channel
		msg := <-broadcast

		// Send message to all clients
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	db, err := sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable password=123456 host=localhost")
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	defer db.Close()

	router := gin.Default()

	userService := service.NewUserService(db)

	router.Use(cors.Default())

	// Route to handle user sign-up
	router.POST("/signup", userService.SignUp)

	// Route to handle user log-in
	router.POST("/login", userService.LogIn)

	router.GET("/ws", handleWebSocket)

	go handleMessages()

	// Start the server
	port := 8080 // You can change the port as needed
	fmt.Printf("Server is running on port %d...\n", port)
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
