package models

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	once = sync.Once{}
	db   *gorm.DB
)

func DB() *gorm.DB {
	once.Do(func() {
		var err error
		if os.Getenv("DATABASE_URL") == "" {
			log.Fatal("Missing database URL! Please set the \"DATABASE_URL\" environment variable.")
		}
		dbUrl, err := url.Parse(os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Fatal("Invalid database URL! Please check the \"DATABASE_URL\" environment variable.")
		}

		username := dbUrl.User.Username()
		password, isSet := dbUrl.User.Password()
		if password == "" || !isSet {
			log.Fatal("Invalid database URL! Please check the \"DATABASE_URL\" environment variable.")
		}
		host, port := dbUrl.Hostname(), dbUrl.Port()
		if port == "" {
			port = "5432" // default Postgres port
		}

		dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s", username, password, dbUrl.Path[1:], host, port)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			AllowGlobalUpdate: false,
			Logger:            logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		// Connection pooling
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get database instance: %v", err)
		}
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)

		//db.AutoMigrate(Club{})
		//db.AutoMigrate(Player{})
	})

	return db
}
