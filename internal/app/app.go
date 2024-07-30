package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/julienschmidt/httprouter"
	"github.com/rlapenok/messagio/internal/handlers"
	"github.com/rlapenok/messagio/internal/state"
	"github.com/rlapenok/messagio/internal/utils"
)

func RunApp() {

	router := httprouter.New()
	router.POST("/send_message", handlers.SendMessage)
	router.GET("/get_stats", handlers.GetSats)
	server := &http.Server{
		Addr:    ":7070",
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Logger.Sugar().Fatalf("ListenAndServeTLS: %v", err)
		}
	}()
	utils.Logger.Info("Server is running on port 7070")

	<-stop
	utils.Logger.Info("Shutting down the server...")

	// Завершение сервера с использованием контекста
	if err := server.Shutdown(context.Background()); err != nil {
		utils.Logger.Sugar().Fatalf("Server Shutdown Failed:%v", err)
	}
	if err := state.State.Close(); err != nil {

		utils.Logger.Sugar().Fatalf("State Shutdown Failed:%v", err)

	}
	utils.Logger.Info("Server gracefully stopped")

}
