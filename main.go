package main

import (
	"bankDepositRating/repos"
	"fmt"
)

func getValidInt() int {
	var i int
	for {
		_, err := fmt.Scanf("%d", &i)
		if err != nil || i < 0 {
			fmt.Println("Некорректный ввод! Введите положительное целое число")
		} else {
			return i
		}
	}

}

func main() {
	bankDeposits := repos.GetBankDeposits()

	fmt.Println("Введите начальную сумму")
	initialSum := getValidInt()

	fmt.Println("Введите сумму ежемесячного пополнения")
	monthlyPayment := getValidInt()

	fmt.Println("Введите срок инвестирования в месяцах")
	investmentPeriod := getValidInt()

	fmt.Printf("Итоговые вложения - %d\n", initialSum+monthlyPayment*investmentPeriod)

	repos.MakeCalculations(bankDeposits, initialSum, monthlyPayment, investmentPeriod)
	repos.SortByEndSum(bankDeposits)
	for _, bankDeposit := range *bankDeposits {
		var capitalization string
		if bankDeposit.Capitalization {
			capitalization = "есть"
		} else {
			capitalization = "нет"
		}
		fmt.Printf("%s %s (срок (мес.) - %d, ставка - %f%%, капитализация - %s) - Выручка: %f, Итог: %f",
			bankDeposit.BankName, bankDeposit.DepositName, bankDeposit.NumberOfMonths, bankDeposit.Rate,
			capitalization, bankDeposit.TotalRevenue, bankDeposit.EndSum)
		fmt.Println()
	}
}
