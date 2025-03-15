package database

import (
	"context"
	"gobank1/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func ConnectDB(DBurl string) (*pgxpool.Pool, error) {
	log := logger.GetLogger()
	log.Info().Str("DBurl", DBurl).Msg("Попытка подключения к базе данных")

	var err error
	Pool, err = pgxpool.New(context.Background(), DBurl)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка подключения к базе данных")
		return nil, err
	}

	err = Pool.Ping(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при проверке подключения к базе данных")
		return nil, err
	}

	log.Info().Msg("Успешное подключение к базе данных")
	return Pool, nil
}

func CloseDB(pool *pgxpool.Pool) {
	log := logger.GetLogger()
	if pool != nil {
		log.Info().Msg("Закрытие соединения с БД")
		pool.Close()
	} else {
		log.Warn().Msg("Попытка закрыть соединение с БД, но оно уже закрыто или не существует")
	}
}
