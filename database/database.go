package database

import (
	"fmt"
	"github.com/MihaPecnik/order-maching-system/config"
	"github.com/MihaPecnik/order-maching-system/models"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase(config *config.Config) (*Database, error) {
	conn, err := gorm.Open(postgres.Open(config.PostgresURI), &gorm.Config{})

	if err != nil {
		fmt.Println("database connection error")
		return nil, err
	}
	fmt.Println("successfully created database connection")

	if config.Migrate {
		err := conn.AutoMigrate(&models.Table{})
		if err != nil {
			fmt.Println(err.Error())
		}
		if config.Populate {
			PopulateDatabase(conn)
		}
	}

	return &Database{db: conn}, nil
}

func PopulateDatabase(db *gorm.DB) {
	orders := []models.Table{
		{
			UserId:   1,
			Value:    200.0,
			Quantity: 2,
			Buy:      true,
			Ticker:   "APPL",
		},
		{
			UserId:   2,
			Value:    200.1,
			Quantity: 3,
			Buy:      true,
			Ticker:   "APPL",
		},
		{
			UserId:   3,
			Value:    200.2,
			Quantity: 1,
			Buy:      true,
			Ticker:   "APPL",
		},
		{
			UserId:   4,
			Value:    200.3,
			Quantity: 5,
			Buy:      true,
			Ticker:   "APPL",
		},
		{
			UserId:   5,
			Value:    201.0,
			Quantity: 2,
			Buy:      false,
			Ticker:   "APPL",
		},
		{
			UserId:   6,
			Value:    201.1,
			Quantity: 3,
			Buy:      false,
			Ticker:   "APPL",
		},
		{
			UserId:   7,
			Value:    201.2,
			Quantity: 1,
			Buy:      false,
			Ticker:   "APPL",
		},
		{
			UserId:   8,
			Value:    201.3,
			Quantity: 5,
			Buy:      false,
			Ticker:   "APPL",
		},
	}
	db.Create(&orders)
}
