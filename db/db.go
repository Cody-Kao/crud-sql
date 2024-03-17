package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Addr     string
	Port     int
	UserName string
	DbName   string
}

type Book struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
	Price  int    `json:"price"`
}

func (Book) TableName() string {
	return "book"
}

func getConfig() (*Config, error) {
	config := &Config{}
	file, err := os.ReadFile("config/config.json")
	if err != nil {
		fmt.Println("Error when reading file")
		return nil, err
	}
	err = json.NewDecoder(strings.NewReader(string(file))).Decode(config)
	if err != nil {
		fmt.Println("Error when decoding json file")
		return nil, err
	}

	return config, nil
}

func ConnectDB() (*sql.DB, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}
	password := os.Getenv("POSTGRES_PASSWORD")
	fmt.Println("password:", password)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Taipei",
		config.Addr, config.Port, config.UserName, password, config.DbName)
	DB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return DB, nil
}
