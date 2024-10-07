package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	MongoDBHost            string
	MongoDBPort            string
	DatabaseName           string
	CollectionNameDetik    string
	CollectionNameLiputan6 string
	DetikMaxPage           int
	LiputanMaxPage         int
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
	DetikMaxPage, err = strconv.Atoi(os.Getenv("DETIK_MAX_PAGE"))
	if err != nil {
		log.Fatalf("Error converting DETIK_MAX_PAGE to int: %v", err)
	}

	LiputanMaxPage, err = strconv.Atoi(os.Getenv("LIPUTAN_MAX_PAGE"))
	if err != nil {
		log.Fatalf("Error converting LIPUTAN_MAX_PAGE to int: %v", err)
	}
}

func init() {
	LoadConfig()
}
