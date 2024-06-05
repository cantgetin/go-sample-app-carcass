package main

import (
	"awesomeProject/config"
	"fmt"
	"github.com/caarlos0/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg := &config.Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("failed to retrieve env variables, %v", err)
	}

	db, err := InitDB(cfg)
	if err != nil {
		log.Fatalf("failed to init db, %v", err)
	}

	var result int
	err = db.Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		log.Fatalf("failed to select 1, %v", err)
	}

	fmt.Printf("Query result: %d\n", result)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err = db.Raw("SELECT 1").Scan(&result).Error
			if err != nil {
				log.Fatalf("failed to select 1, %v", err)
			}
			fmt.Println("Selected 1")

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Check the console for the message."))
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// Create the DSN string using the configuration values
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbHost, cfg.DbPort, cfg.DbUsername, cfg.DbPassword, cfg.DbDatabase,
	)

	newLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormLogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormLogger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("can't connect to pg instance, %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.PgIdleConn)
	sqlDB.SetMaxOpenConns(cfg.PgMaxOpenConn)

	go func() {
		t := time.NewTicker(cfg.PgPingInterval)
		defer t.Stop()

		for range t.C {
			if err := sqlDB.Ping(); err != nil {
				log.Printf("error pinging database: %v", err)
			}
		}
	}()

	return db, nil
}
