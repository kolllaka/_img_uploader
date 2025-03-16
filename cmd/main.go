package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kolllaka/_img_uploader/internal/img"
	"github.com/kolllaka/_img_uploader/internal/model"
	"github.com/kolllaka/_img_uploader/internal/router"
	"github.com/kolllaka/_img_uploader/internal/service"
	storage "github.com/kolllaka/_img_uploader/internal/storage/sqlite"
	"github.com/kolllaka/_img_uploader/pkg/config"
	"github.com/kolllaka/_img_uploader/pkg/sqlite"
)

const (
	config_path = "./config.yaml"
)

func main() {
	cfg := model.NewConfig()
	config.MustLoadByPath(config_path, cfg)

	db, err := sqlite.NewDB(cfg.DB.Path)
	if err != nil {
		panic(err)
	}

	store := storage.NewImgStore(cfg, db)

	transImg := img.NewImage()

	service := service.New(cfg, store, transImg)

	go func() {
		clearFiles(service)
		ticker := time.NewTicker(cfg.Files.ClearDelay * time.Hour)
		for range ticker.C {
			clearFiles(service)
		}
	}()

	server := router.New(service)
	mux := server.Init()

	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads/"))))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	fmt.Println("Starting server")

	http.ListenAndServe(":8080", mux)
}

func clearFiles(service service.Service) {
	cherr := make(chan error)
	defer close(cherr)
	go func() {
		for err := range cherr {
			fmt.Printf("Error deleting expired files: %v\n", err)
		}
	}()

	count := service.DeleteExpiresFiles(cherr)

	fmt.Printf("Deleted %d expired files %s\n", count, time.Now().Format(time.RFC822))
}
