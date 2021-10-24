package config

import (
	"log"
	"os"

	controllers "example.com/gorestapi/controllers"
	"github.com/go-pg/pg/v9"
	"github.com/magiconair/properties"
)

// Connecting to db
func ConnectDBandProperties() *pg.DB {
	opts := &pg.Options{
		User:     "postgres",
		Password: "admin",
		Addr:     "localhost:5432",
		Database: "mydb",
	}

	var db *pg.DB = pg.Connect(opts)
	if db == nil {
		log.Printf("Failed to connect")
		os.Exit(100)
	}
	log.Printf("Connected to db")
	controllers.InitiateDB(db)
	p := properties.MustLoadFile("configs/properties/config.properties", properties.UTF8)
	controllers.InitiateProp(p)

	return db
}
