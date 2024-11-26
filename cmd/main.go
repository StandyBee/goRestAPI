package main

import (
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	serv "gorestAPI"
	"gorestAPI/pkg/handler"
	"gorestAPI/pkg/repository"
	"gorestAPI/pkg/service"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if err := initConfig(); err != nil {
		logrus.Fatalf("Error reading config file, %s", err)
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatal("Error loading .env file")
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		Username: viper.GetString("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: viper.GetString("DATABASE"),
		SSLMode:  viper.GetString("DB_SSLMODE"),
	})
	if err != nil {
		logrus.Fatalf("Error connecting to database, %s", err)
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(serv.Server)

	go func() {
		err = srv.Run(viper.GetString("port"), handlers.InitRoutes())
		if err != nil {
			logrus.Fatalf("error occured while running server: %s", err.Error())
		}
	}()

	logrus.Print("TODO app started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Shutting down server...")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("Error occured while shutting down server, %s", err)
	}

	if err := db.Close(); err != nil {
		logrus.Fatalf("Error occured while closing database, %s", err)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
