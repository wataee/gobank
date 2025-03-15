package checkerName

import (
	"bufio"
	"context"
	"fmt"
	"time"

	"gobank1/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunCheckerName(scanner *bufio.Scanner, pool *pgxpool.Pool) string {
	log := logger.GetLogger()
	var name string
	for {
		fmt.Println("Введите имя:")
		scanner.Scan()
		name = scanner.Text()

		var exists bool
		err := pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE name = $1)", name).Scan(&exists)
		if err != nil {
			log.Error().Err(err).Msg("Ошибка при проверке имени в базе данных")
			return ""
		}

		if !exists {
			log.Warn().Msg("Пользователь не найден в БД")
			fmt.Println("Пользователь не найден в БД")
		} else {
			log.Info().Str("name", name).Msg("Пользователь найден в БД")
			fmt.Println("Пользователь найден в БД")
			time.Sleep(1 * time.Second)
			return name
		}
	}
}

func RunCheckerName2(scanner *bufio.Scanner, pool *pgxpool.Pool) string {
	log := logger.GetLogger()
	var name string
	for {
		scanner.Scan()
		name = scanner.Text()

		var exists bool
		err := pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE name = $1)", name).Scan(&exists)
		if err != nil {
			log.Error().Err(err).Msg("Ошибка при проверке имени в базе данных")
			return ""
		}

		if !exists {
			log.Warn().Msg("Пользователь не найден в БД. Повторите попытку")
			fmt.Println("Пользователь не найден в БД. Повторите попытку")
		} else {
			log.Info().Str("name", name).Msg("Пользователь найден")
			fmt.Println("Пользователь найден")
			time.Sleep(1 * time.Second)
			return name
		}
	}
}

func RunCheckerName3(scanner *bufio.Scanner, pool *pgxpool.Pool) string {
	log := logger.GetLogger()
	var name string
	for {
		fmt.Println("Введите имя:")
		scanner.Scan()
		name = scanner.Text()

		var exists bool
		err := pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE name = $1)", name).Scan(&exists)
		if err != nil {
			log.Error().Err(err).Msg("Ошибка при проверке имени в базе данных")
			return ""
		}

		if !exists {
			log.Info().Str("name", name).Msg("Пользователь не найден, принимаем имя")
			return name
		} else {
			log.Warn().Str("name", name).Msg("Пользователь уже есть в БД. Повторите попытку")
			fmt.Println("Пользователь уже есть в БД. Повторите попытку")
			time.Sleep(1 * time.Second)
		}
	}
}
