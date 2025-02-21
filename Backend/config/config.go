package config

import (
	"log"
	"os"
)

type Config struct {
	Port       string
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	JWTSecret  string
	OpenAIKey  string
}

var ConfigInstance Config

func Init() {
	ConfigInstance = Config{
		Port:       os.Getenv("PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		OpenAIKey:  os.Getenv("OPENAI_KEY"),
	}

	if ConfigInstance.Port == "" {
		ConfigInstance.Port = "8080" // Default port
		log.Println("Using default port 8080")
	}

	// Initialize database connection
	InitDB()
}