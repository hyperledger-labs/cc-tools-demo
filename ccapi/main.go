package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/goledgerdev/ccapi/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Create gin handler and start server
	r := gin.Default()
	go server.Serve(r, ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	cancel()
}
