package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Wait for interrupt signal to gracefully shut down the server with a timeout of 5 seconds.

// kill (no param) default send syscall.SIGTERM
// kill -2 is syscall.SIGINT
// kill -9 is syscall. SIGKILL but can"t be caught, so don't need to add it
func (srv *HTTPServer) waitOSInterruptSignal() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
}

func (srv *HTTPServer) terminateServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.Server.Shutdown(ctx)
	if err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}

	log.Println("Server exiting")
}

func (srv *HTTPServer) gracefulShutDown() {
	srv.waitOSInterruptSignal()
	srv.Domain.GetQueueService().Stop()
	srv.terminateServer()
}
