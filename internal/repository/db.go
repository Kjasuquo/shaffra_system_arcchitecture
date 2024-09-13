package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"shaffra_assessment/config"
	"shaffra_assessment/internal/models"
)

type PostgresDB struct {
	instance *gorm.DB
	config   *config.Config
}

// Connect is used to connect the database instance
func (p *PostgresDB) Connect(config config.Config) error {

	var dsn string
	databaseUrl := config.DatabaseUrl

	if databaseUrl == "" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
	} else {
		dsn = databaseUrl
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	p.instance = db
	p.config = &config

	err = p.instance.AutoMigrate(&models.User{})

	if err != nil {
		log.Println("error migrating tables: ", err)
		log.Panic(err)
	}

	log.Println("Database Connected Successfully")
	return nil
}

func (p *PostgresDB) CloseConnection() error {
	if p.instance != nil {
		connection, err := p.instance.DB()
		if err != nil {
			return err
		}

		connection.Close()
	}
	return nil
}

func (p *PostgresDB) RestartConnection(config config.Config) error {
	if p.instance != nil {
		p.CloseConnection()
	}

	return p.Connect(config)
}
