package cli

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"gobank1/internal/account"
	"gobank1/internal/transactions"
	"gobank1/pkg/database"
)

func RunCLI() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("-------------------------------------------")
		fmt.Println("Выберите действие: \n \n 1. Зарегистрироваться \n 2. Статистика \n 3. Снять деньги со счёта \n 4. Пополнение счёта \n 5. История транзакций \n 6. Перевод другому пользователю \n 7. Выход")
		fmt.Println("-------------------------------------------")
		scanner.Scan()
		switch scanner.Text() {
		case "1":
			account.RunRegistration(database.Pool)
				time.Sleep(1 * time.Second)
		case "2":
			account.RunStats(database.Pool)
				time.Sleep(1 * time.Second)
		case "3":
			transactions.ProcessTransaction(database.Pool, "withdrawal")
				time.Sleep(1 * time.Second)
		case "4":
			transactions.ProcessTransaction(database.Pool, "deposit")
				time.Sleep(1 * time.Second)
		case "5":
			account.RunTrHistory(database.Pool)
				time.Sleep(1 * time.Second)
		case "6":
			transactions.RunTransfer(database.Pool)
				time.Sleep(1 * time.Second)
		case "7":
			fmt.Println("Выход из программы")
			return
		default:
			fmt.Println("Некорректный ввод. Попробуйте снова")
			time.Sleep(1 * time.Second)
		}
	}
}