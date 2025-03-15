package transactions

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"

	"gobank1/pkg/checkerName"
	"gobank1/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunTransfer(pool *pgxpool.Pool) {
	var nameFrom, nameTo string
	var amount float64

	scanner := bufio.NewScanner(os.Stdin)

	log := logger.GetLogger()

	log.Info().Msg("Запуск перевода")

	fmt.Println("Введите имя отправителя:")
	nameFrom = checkerName.RunCheckerName2(scanner, pool)

	fmt.Println("Введите имя получателя:")
	nameTo = checkerName.RunCheckerName2(scanner, pool)

	fmt.Println("Введите сумму:")
	scanner.Scan()
	amount, _ = strconv.ParseFloat(scanner.Text(), 64)

	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка начала транзакции")
		fmt.Println("Ошибка начала транзакции")
		return
	}
	defer tx.Rollback(ctx)

	var balanceFrom float64
	err = tx.QueryRow(ctx, "SELECT balance FROM users WHERE name = $1", nameFrom).Scan(&balanceFrom)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка получения баланса отправителя")
		fmt.Println("Ошибка получения баланса отправителя:", err)
		return
	}

	if balanceFrom < amount {
		log.Warn().Msg("Недостаточно средств для перевода")
		fmt.Println("Ошибка: недостаточно средств для перевода")
		return
	}

	_, err = tx.Exec(ctx, "UPDATE users SET balance = balance - $1 WHERE name = $2 AND balance >= $1", amount, nameFrom)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при списании средств у отправителя")
		fmt.Println("Ошибка при списании средств у отправителя:", err)
		return
	}

	_, err = tx.Exec(ctx, "UPDATE users SET balance = balance + $1 WHERE name = $2", amount, nameTo)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при зачислении средств получателю")
		fmt.Println("Ошибка при зачислении средств получателю:", err)
		return
	}

	_, err = tx.Exec(ctx,
		"INSERT INTO trhistory (user_id, amount, transaction_type, created_at) "+
			"SELECT id, $1, 'debit', NOW() FROM users WHERE name = $2",
		amount, nameFrom)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при записи транзакции отправителя")
		fmt.Println("Ошибка при записи транзакции отправителя:", err)
		return
	}

	_, err = tx.Exec(ctx,
		"INSERT INTO trhistory (user_id, amount, transaction_type, created_at) "+
			"SELECT id, $1, 'credit', NOW() FROM users WHERE name = $2",
		amount, nameTo)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при записи транзакции получателя")
		fmt.Println("Ошибка при записи транзакции получателя:", err)
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при коммите транзакции")
		fmt.Println("Ошибка при коммите транзакции:", err)
		return
	}

	log.Info().Msg("Перевод успешно выполнен!")
	fmt.Println("Перевод успешно выполнен!")
}
