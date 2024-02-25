package main

import (
	"database/sql"
	"fmt"

	service "meetcode-backend/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

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

	// Start the server
	port := 8080 // You can change the port as needed
	fmt.Printf("Server is running on port %d...\n", port)
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
