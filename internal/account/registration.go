package account

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gobank1/internal/models"
	"gobank1/pkg/checkerName"
	"gobank1/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunRegistration(pool *pgxpool.Pool) {
	log := logger.GetLogger()

	var user models.User
	scanner := bufio.NewScanner(os.Stdin)
	log.Info().Msg("Запуск регистрации пользователя")

	name := checkerName.RunCheckerName3(scanner, pool)
	log.Debug().Str("name", name).Msg("Получено имя пользователя")

	if strings.TrimSpace(name) == "" {
		log.Error().Msg("Ошибка: имя не может быть пустым")
		return
	} else if _, err := strconv.Atoi(name); err == nil {
		log.Error().Msg("Ошибка: имя не может состоять только из цифр")
		return
	}

	fmt.Print("Введите начальный баланс: ")
	scanner.Scan()
	input := scanner.Text()
	log.Debug().Str("input", input).Msg("Получен ввод баланса")

	balance, err := strconv.ParseFloat(strings.Replace(input, ",", ".", 1), 64)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка: введите корректное число для баланса")
		return
	}

	var userid int
	var createdAt time.Time
	query := `INSERT INTO users (name, balance) VALUES ($1, $2) RETURNING id, created_at`
	log.Debug().Str("query", query).Msg("Выполняется запрос к БД")

	err2 := pool.QueryRow(context.Background(), query, name, balance).Scan(&userid, &createdAt)
	if err2 != nil {
		log.Error().Err(err2).Msg("Ошибка при добавлении данных пользователя в БД")
		return
	}

	user.Name = name
	user.Balance = balance
	user.ID = userid
	user.CreatedAt = createdAt

	log.Info().Int("userID", user.ID).Str("name", user.Name).Float64("balance", user.Balance).Msg("Пользователь успешно зарегистрирован")

	fmt.Println("-------------------------------------------")
	fmt.Printf("Пользователь зарегистрирован.\nИмя: %s\nБаланс: %.2f\nДата создания: %s\n", user.Name, user.Balance, user.CreatedAt)
	fmt.Println("-------------------------------------------")
}
