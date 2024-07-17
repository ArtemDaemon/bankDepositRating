package repos

import (
	"bankDepositRating/database"
)

type BankDeposit struct {
	Id             int
	BankName       string
	DepositName    string
	Rate           float64
	NumberOfMonths int
}

func GetBankDeposits() *[]BankDeposit {
	rows := database.ExecQuery("SELECT * FROM bank_deposits;")
	defer rows.Close()
	var bankDeposits []BankDeposit
	for rows.Next() {
		var bankDeposit BankDeposit
		err := rows.Scan(&bankDeposit.Id, &bankDeposit.BankName, &bankDeposit.DepositName, &bankDeposit.Rate,
			&bankDeposit.NumberOfMonths)
		if err != nil {
			panic(err)
		}
		bankDeposits = append(bankDeposits, bankDeposit)
	}
	return &bankDeposits
}
