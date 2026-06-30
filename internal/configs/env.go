package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/alexedwards/argon2id"
	"github.com/joho/godotenv"
)

type env struct {
	PORT            string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBHost          string
	DBSchema        string
	JwtSecret       string
	ServerENV       string
	DefaultUser     string
	DefaultEmail    string
	DefaultPassword string
	N8NBaseURL      string
	R2AccountID     string
	R2AccessKeyID   string
	R2SecretKeyID   string
	R2Token         string
	R2PublicDomain  string
	R2Bucket        string
	SMTPHost        string
	SMTPPort        int
	SMTPName        string
	SMTPEmail       string
	SMTPPassword    string
}

var Env *env

func SetupENV() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err.Error())
	}

	defaultPassword, _ := argon2id.CreateHash(os.Getenv("DEFAULT_PASSWORD"), argon2id.DefaultParams)

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		smtpPort = 587
	}

	Env = &env{
		PORT:            os.Getenv("PORT"),
		DBPort:          os.Getenv("DB_PORT"),
		DBUser:          os.Getenv("DB_USER"),
		DBPassword:      os.Getenv("DB_PASSWORD"),
		DBName:          os.Getenv("DB_NAME"),
		DBHost:          os.Getenv("DB_HOST"),
		DBSchema:        os.Getenv("DB_SCHEMA"),
		JwtSecret:       os.Getenv("JWT_SECRET"),
		ServerENV:       os.Getenv("SERVER_ENV"),
		DefaultUser:     os.Getenv("DEFAULT_USER"),
		DefaultEmail:    os.Getenv("DEFAULT_EMAIL"),
		DefaultPassword: defaultPassword,
		N8NBaseURL:      os.Getenv("N8N_BASE_URL"),
		R2AccountID:     os.Getenv("R2_ACCOUNT_ID"),
		R2AccessKeyID:   os.Getenv("R2_ACCESS_KEY_ID"),
		R2SecretKeyID:   os.Getenv("R2_SECRET_KEY_ID"),
		R2Token:         os.Getenv("R2_TOKEN"),
		R2PublicDomain:  os.Getenv("R2_PUBLIC_DOMAIN"),
		R2Bucket:        os.Getenv("R2_BUCKET"),
		SMTPHost:        os.Getenv("SMTP_HOST"),
		SMTPPort:        smtpPort,
		SMTPName:        os.Getenv("SMTP_NAME"),
		SMTPEmail:       os.Getenv("SMTP_EMAIL"),
		SMTPPassword:    os.Getenv("SMTP_PASSWORD"),
	}
}
