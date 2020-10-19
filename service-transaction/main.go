package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rifarihadts/microservice/service-product/config"
	"github.com/rifarihadts/microservice/service-product/database"
	"github.com/rifarihadts/microservice/service-product/handler"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Panic(err)
		return
	}

	db, err := initDB(cfg.Database)

	router := mux.NewRouter()
	authMiddleware := handler.AuthMiddleware{
		AuthService: cfg.AuthService,
	}
	transactionHandler := handler.Transaction{Db: db}

	router.Handle("/transaction/add", authMiddleware.ValidateAuth(transactionHandler.AddTransaction))
	router.Handle("/transactions", authMiddleware.ValidateAuth(transactionHandler.GetTransaction))

	fmt.Printf("Server listen on :%s", cfg.Port)
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), router))
}

func getConfig() (config.Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	viper.SetConfigName("config.yml")

	if err := viper.ReadInConfig(); err != nil {
		return config.Config{}, err
	}

	var cfg config.Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return config.Config{}, err
	}

	return cfg, nil
}

func initDB(dbConfig config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DbName, dbConfig.Config)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&database.Transaction{})
	if err != nil {
		return nil, err
	}

	return db, nil
}