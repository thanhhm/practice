package main

import (
	"log"
	"os"

	"github.com/subosito/gotenv"
)

type config struct {
	ginPort              string
	oracleDataSourceName string
}

func newConfig() *config {
	if err := gotenv.Load(); err != nil {
		log.Fatal("Missing environment file")
	}

	return &config{
		ginPort:              os.Getenv("GIN_PORT"),
		oracleDataSourceName: os.Getenv("ORACLE_DATA_SOURCE_NAME"),
	}
}
