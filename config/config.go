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

func getEnvAsInt(key string) int {
	valueStr := os.Getenv(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("Error converting %s to int: %v", key, err)
	}
	return value
}

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
	DetikMaxPage = getEnvAsInt("DETIK_MAX_PAGE")
	LiputanMaxPage = getEnvAsInt("LIPUTAN_MAX_PAGE")
}

func init() {
	LoadConfig()
}
