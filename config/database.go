package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"os"
)

// GetConnection get connection database
func GetConnection() *gorm.DB {
	c := GetConfig()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.Database.User, c.Database.Pass, c.Database.Server, c.Database.Port, c.Database.Database)
	}

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
