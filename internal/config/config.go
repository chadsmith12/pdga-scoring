package config

import (
	"fmt"
	"os"
	"strings"
)

type DatabaseConfig struct {
    Host string
    Username string
    Password string
    Port string
    DatabaseName string
    SSLMode bool
}

func (dbConfig DatabaseConfig) String() string {
    return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
        dbConfig.Username,
        dbConfig.Password,
        dbConfig.Host,
        dbConfig.Port,
        dbConfig.DatabaseName)
}

func LoadDatabase() (DatabaseConfig, error) {
    dbConfig := DatabaseConfig{}

    host, ok := os.LookupEnv("PG_HOST")
    if !ok {
        return dbConfig, fmt.Errorf("no PG_HOST env variable set")
    }
    dbConfig.Host = host

    username, ok := os.LookupEnv("PG_USERNAME")
    if !ok {
        return dbConfig, fmt.Errorf("no PG_USERNAME env variable set")
    }
    dbConfig.Username = username

    password, ok := os.LookupEnv("PG_PASSWORD")
    if !ok {
        return dbConfig, fmt.Errorf("no PG_PASSWORD env variable set")
    }
    dbConfig.Password = password 

    port, ok := os.LookupEnv("PG_PORT")
    if !ok {
        return dbConfig, fmt.Errorf("no PG_PORT env variable set")
    }
    dbConfig.Port = port 

    dbName, ok := os.LookupEnv("PG_DBNAME")
    if !ok {
        return dbConfig, fmt.Errorf("no PG_DBNAME env variable set")
    }
    dbConfig.DatabaseName = dbName

    dbConfig.SSLMode = loadSSlMode()

    return dbConfig, nil
}

func loadSSlMode() bool {
    sslMode, ok := os.LookupEnv("PG_SSLMODE")
    if !ok {
        return false
    }

    if strings.ToLower(sslMode) == "true" {
        return true
    }

    return false
}
