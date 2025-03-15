package account

import (
	"bufio"
	"context"
	"os"
	"fmt"

	"gobank1/internal/models"
	"gobank1/pkg/checkerName"
	"gobank1/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunStats(pool *pgxpool.Pool) {
	log := logger.GetLogger()
	log.Info().Msg("Запуск RunStats")

	var name string
	scanner := bufio.NewScanner(os.Stdin)
	name = checkerName.RunCheckerName(scanner, pool)
	log.Info().Str("user_name", name).Msg("Получено имя пользователя")

	rows, err := pool.Query(context.Background(), "SELECT * FROM users WHERE name = $1", name)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка при выполнении запроса")
		return
	}
	defer rows.Close()

	var user models.User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Balance, &user.CreatedAt); err != nil {
			log.Error().Err(err).Msg("Ошибка при сканировании данных")
			return
		}
	}

	log.Info().Interface("user", user).Msg("Получена статистика пользователя")

	log.Info().Msg("Вывод статистики пользователя")
	fmt.Println("-------------------------------------------")
	fmt.Printf("Статистика пользователя %s\n\nID: %v\nИмя: %v\nБаланс: %v\nДата регистрации: %v\n", user.Name, user.ID, user.Name, user.Balance, user.CreatedAt)
	fmt.Println("-------------------------------------------")
}

func RunCheckerName(name string, pool *pgxpool.Pool) {
	panic("unimplemented")
}
