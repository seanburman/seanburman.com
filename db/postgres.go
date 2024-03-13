package db

import (
	"log"

	"github.com/seanburman/seanburman.com/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDriver() *gorm.DB {
	client, err := gorm.Open(postgres.Open(config.Env().DATABASE_URL), &gorm.Config{})
	if err != nil {
		log.Panicln(err.Error())
	} else {
		log.Println("Connected to Postres DB...")
	}
	return client
}
