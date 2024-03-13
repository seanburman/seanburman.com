package db

import "gorm.io/gorm"

var Instance *Database

type Database struct {
	Postgres *gorm.DB
}
