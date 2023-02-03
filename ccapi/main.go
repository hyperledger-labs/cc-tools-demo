package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/goledgerdev/ccapi/chaincode"
	"github.com/goledgerdev/ccapi/server"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Create gin handler and start server
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:8080", // Test addresses
			"*",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Origin", "Content-Type"},
		AllowCredentials: true,
	}))
	go server.Serve(r, ctx)

	// Register to chaincode events
	go chaincode.Event(os.Getenv("CHANNEL"), os.Getenv("CCNAME"), "eventName", func(ccEvent *fab.CCEvent) {
		log.Println("Received CC event: ", ccEvent)
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	cancel()
}
