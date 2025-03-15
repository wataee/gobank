package transactions

import (
	"bufio"
	"context"
	"os"
	"strconv"
	"time"
	"fmt"

	"gobank1/internal/models"
	"gobank1/pkg/checkerName"
	"gobank1/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ProcessTransaction(pool *pgxpool.Pool, transactionType string) {
	log := logger.GetLogger()
	log.Info().Str("transaction_type", transactionType).Msg("Запуск ProcessTransaction")

	scanner := bufio.NewScanner(os.Stdin)
	name := checkerName.RunCheckerName(scanner, pool)
	log.Info().Str("user_name", name).Msg("Получено имя пользователя")

	var user models.User
	var amount float64

	fmt.Printf("Введите сумму: ")
	scanner.Scan()
	amount, err := strconv.ParseFloat(scanner.Text(), 64)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка парсинга суммы")
		return
	}

	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при начале транзакции")
		return
	}
	defer tx.Rollback(ctx)

	err = pool.QueryRow(ctx, "SELECT id, balance FROM users WHERE name = $1", name).Scan(&user.ID, &user.Balance)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при получении данных пользователя")
		return
	}

	log.Info().Int("user_id", user.ID).Float64("balance", user.Balance).Msg("Данные пользователя загружены")

	if transactionType == "withdrawal" && user.Balance < amount {
		log.Warn().Msg("Недостаточно средств на балансе")
		return
	}

	updateQuery := "UPDATE users SET balance = balance + $1 WHERE id = $2"
	if transactionType == "withdrawal" {
		updateQuery = "UPDATE users SET balance = balance - $1 WHERE id = $2"
	}

	_, err = tx.Exec(ctx, updateQuery, amount, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при обновлении баланса")
		return
	}

	transaction := models.TransactionHistory{
		UserId:          user.ID,
		Amount:          amount,
		TransactionType: transactionType,
		CreatedAt:       time.Now(),
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO trHistory (user_id, amount, transaction_type, created_at) 
		VALUES ($1, $2, $3, $4)`,
		transaction.UserId, transaction.Amount, transaction.TransactionType, transaction.CreatedAt)

	if err != nil {
		log.Error().Err(err).Msg("Ошибка при записи транзакции")
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при коммите транзакции")
		return
	}

	log.Info().Str("transaction_type", transactionType).Str("user_name", name).Float64("amount", amount).Msg("Транзакция успешно выполнена")
	fmt.Printf("%s | успешно выполнено. Текущий баланс %s: %.2f\n", transactionType, name, user.Balance+amount)
}