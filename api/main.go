package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/0x726f6f6b6965/task/internal/config"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

func main() {
	godotenv.Load()
	path := os.Getenv("CONFIG")
	var cfg config.Config
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("read yaml error", err)
		return
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal("unmarshal yaml error", err)
		return
	}

	mux, cleanup, err := initApplication(context.Background(), &cfg)
	if err != nil {
		log.Fatal("initialize application error", err)
	}
	defer cleanup()

	log.Printf("server listening; port: %d", cfg.Rest.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Rest.Port), mux); err != nil {
		log.Fatalf("failed to serve; err: %v", err)
		return
	}
}
