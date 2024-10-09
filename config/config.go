package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	MongoDBHost            string
	MongoDBPort            string
	DatabaseName           string
	CollectionNameDetik    string
	CollectionNameLiputan6 string
)

func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	MongoDBHost = os.Getenv("MONGO_DB_HOST")
	MongoDBPort = os.Getenv("MONGO_DB_PORT")
	DatabaseName = os.Getenv("DATABASE_NAME")
	CollectionNameDetik = os.Getenv("COLLECTION_NAME_DETIK")
	CollectionNameLiputan6 = os.Getenv("COLLECTION_NAME_LIPUTAN")
}

func init() {
	LoadConfig()
}
