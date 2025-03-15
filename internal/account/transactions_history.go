package account

import (
	"bufio"
	"context"
	"os"
	"time"
	"fmt"

	"gobank1/pkg/checkerName"
	"gobank1/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunTrHistory(pool *pgxpool.Pool) {
	log := logger.GetLogger()
	log.Info().Msg("Запуск RunTrHistory")

	scanner := bufio.NewScanner(os.Stdin)
	name := checkerName.RunCheckerName(scanner, pool)
	log.Info().Str("user_name", name).Msg("Получено имя пользователя")

	rows, err := pool.Query(context.Background(), `
		SELECT trhistory.id, trhistory.user_id, trhistory.amount, 
		       trhistory.transaction_type, trhistory.created_at, users.name 
		FROM trhistory 
		JOIN users ON users.id = trhistory.user_id 
		WHERE users.name = $1`, name)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка запроса")
		return
	}
	defer rows.Close()

	var transactions []struct {
		ID             int
		UserID         int
		Amount         float64
		TransactionType string
		CreatedAt      time.Time
		UserName       string
	}

	for rows.Next() {
		var tr struct {
			ID             int
			UserID         int
			Amount         float64
			TransactionType string
			CreatedAt      time.Time
			UserName       string
		}

		err := rows.Scan(&tr.ID, &tr.UserID, &tr.Amount, &tr.TransactionType, &tr.CreatedAt, &tr.UserName)
		if err != nil {
			log.Error().Err(err).Msg("Ошибка чтения строки")
			return
		}

		transactions = append(transactions, tr)
	}

	if len(transactions) == 0 {
		log.Warn().Msg("У пользователя нет транзакций")
		return
	}

	log.Info().Int("transaction_count", len(transactions)).Msg("Вывод истории транзакций")
	fmt.Println("-------------------------------------------")
	fmt.Printf("История транзакций %s:\n", name)
	for _, tr := range transactions {
		fmt.Printf("ID транзакции: %v | ID пользователя: %v | Сумма: %.2f | Тип: %s | Дата: %v\n",
			tr.ID, tr.UserID, tr.Amount, tr.TransactionType, tr.CreatedAt)
	}
	fmt.Println("-------------------------------------------")
}
