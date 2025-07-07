package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/CRobinDev/karsa/config/env"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var (
	db   *sqlx.DB
	once sync.Once
)

func NewPostgresPool(logger *logrus.Logger) *sqlx.DB {
	once.Do(func() {
		dataSourceName := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			env.GetEnv().DBHost, 
			env.GetEnv().DBPort, 
			env.GetEnv().DBUser, 
			env.GetEnv().DBPassword, 
			env.GetEnv().DBName,
		)

		pool, err := sqlx.Connect("pgx", dataSourceName)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"error": err.Error(),
				"host":  env.GetEnv().DBHost,
			}).Errorf("[Database][NewPostgresPool] failed to connect to database 😊")
		}
		pool.SetMaxOpenConns(100)
		pool.SetMaxIdleConns(10)
		pool.SetConnMaxLifetime(60 * time.Minute)

		db = pool
	})

	return db
}
