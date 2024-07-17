package main

import (
	"bankDepositRating/repos"
	"fmt"
)

func main() {
	bankDeposits := repos.GetBankDeposits()
	for _, bankDeposit := range *bankDeposits {
		fmt.Println(bankDeposit.Id, bankDeposit.BankName, bankDeposit.DepositName, bankDeposit.Rate,
			bankDeposit.NumberOfMonths)
	}
}
