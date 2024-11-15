package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var tables = []interface{}{&Owner{}, &Evidence{}}

func init() {
	envErr := godotenv.Load()
	if envErr != nil {
		panic("Error loading .env file")
	}

	var dbErr error
	db, dbErr = gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if dbErr != nil {
		log.Fatalf("Failed to connect to database: %v", dbErr)
	}

	migErr := db.AutoMigrate(tables...)
	if migErr != nil {
		log.Fatalf("Tables migration failing: %v", migErr)
	}

	fmt.Println("Database connection established & tables migrated")
}

type Owner struct {
	PubAddress      string `gorm:"primaryKey"`
	Name            string
	AccessTimestamp *int
	CreationTx      string
	AccessTx        *string
	MasterId        *string
	Master          *Owner `gorm:"foreignKey:MasterId"`
}

type Evidence struct {
	EvId       string `gorm:"primaryKey"`
	Extension  string
	Index      int    // Index of the evidence in the list on blockchain (-1 Initially)
	CreationTx string // 0x0 Initially
	OwnerAddr  string
	Owner      *Owner `gorm:"foreignKey:OwnerAddr"`
}
