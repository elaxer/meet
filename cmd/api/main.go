package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"meet/internal/config"
	"meet/internal/pkg/app/helper"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	_, b, _, _ = runtime.Caller(0)
	rootDir, _ = filepath.Abs(filepath.Dir(b) + "/../..")
)

var db *sql.DB
var cfg *config.Config

func init() {
	err := godotenv.Load(rootDir + "/.env")
	if err != nil {
		panic(err)
	}

	cfg = config.NewConfig(rootDir)

	db, err = helper.LoadDB(cfg.DBConfig)
	if err != nil {
		panic(err)
	}
}

func main() {
	http.Handle("/", httpHandler(db))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerConfig.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	go func() {
		err := server.ListenAndServe()
		log.Println(err)
	}()

	log.Printf("Сервер запущен на адресе %s", server.Addr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Остановка сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
