package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/vanguard"
	library_v1 "github.com/joematpal/library_rpc/internal/library/v1"
	"github.com/joematpal/library_rpc/pkg/library/v1/library_v1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {

	libraryService, err := library_v1.NewLibraryService()
	if err != nil {
		log.Fatalf("new library service: %v", err)
	}
	mux := http.NewServeMux()

	libraryServicePath, libraryServiceHandler := library_v1connect.NewLibraryServiceHandler(libraryService)

	mux.Handle(libraryServicePath, libraryServiceHandler)
	vangaurdLibraryService := vanguard.NewService(
		libraryServicePath,
		libraryServiceHandler,
		vanguard.WithTargetProtocols(vanguard.ProtocolConnect),
		vanguard.WithTargetCodecs(vanguard.CodecJSON),
	)

	transcoder, err := vanguard.NewTranscoder([]*vanguard.Service{
		vangaurdLibraryService,
	})
	if err != nil {
		log.Fatalf("new vanguard: %v", err)
	}

	mux.Handle("/api", http.StripPrefix("/api", transcoder))
	corshandler := addCORS(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(corshandler, &http2.Server{}),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down service....")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func addCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Method", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version, Connect-Timeout-Ms")
		w.Header().Set("Access-Control-Allow-Expose-Headers", "Connect-Protocol-Version, Connect-Timeout-Ms")

		handler.ServeHTTP(w, r)
	})
}
