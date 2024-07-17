package repos

import (
	"bankDepositRating/database"
	"math"
	"sort"
)

type Calculation struct {
	StartMonth int
	EndMonth   int
	StartSum   float64
	Revenue    float64
	Savings    float64
	EndSum     float64
}

type BankDeposit struct {
	Id             int
	BankName       string
	DepositName    string
	Rate           float64
	NumberOfMonths int
	Capitalization bool
	Expandable     bool
	Calculations   []Calculation
	EndSum         float64
	TotalRevenue   float64
}

func GetBankDeposits() *[]BankDeposit {
	rows := database.ExecQuery("SELECT * FROM bank_deposits;")
	defer rows.Close()
	var bankDeposits []BankDeposit
	for rows.Next() {
		var bankDeposit BankDeposit
		err := rows.Scan(&bankDeposit.Id, &bankDeposit.BankName, &bankDeposit.DepositName, &bankDeposit.Rate,
			&bankDeposit.NumberOfMonths, &bankDeposit.Capitalization, &bankDeposit.Expandable)
		if err != nil {
			panic(err)
		}
		bankDeposits = append(bankDeposits, bankDeposit)
	}
	return &bankDeposits
}

func MakeCalculations(bankDeposits *[]BankDeposit, initialSum int, monthlyPayment int, investmentPeriod int) {
	for i := range *bankDeposits {
		bankDeposit := &(*bankDeposits)[i]
		if bankDeposit.Calculations == nil {
			bankDeposit.Calculations = make([]Calculation, 0)
		}
		startMonth := 0
		startSum := float64(initialSum)
		var endSum float64
		for startMonth+bankDeposit.NumberOfMonths <= investmentPeriod {
			var revenue float64
			endMonth := startMonth + bankDeposit.NumberOfMonths
			savings := float64(bankDeposit.NumberOfMonths * monthlyPayment)
			depositLength := float64(31 * bankDeposit.NumberOfMonths)

			newCalculation := Calculation{
				StartMonth: startMonth,
				EndMonth:   endMonth,
				StartSum:   startSum,
				Savings:    savings,
			}

			if bankDeposit.Expandable {
				startSum += savings
				if bankDeposit.Capitalization {

				} else {
					revenue = (startSum * bankDeposit.Rate * (depositLength / 366)) / 100
				}
				endSum = startSum + revenue
			} else {
				if bankDeposit.Capitalization {
					revenue = startSum*math.Pow(1+bankDeposit.Rate/100/12, float64(bankDeposit.NumberOfMonths)) -
						startSum
				} else {
					revenue = (startSum * bankDeposit.Rate * (depositLength / 366)) / 100
				}
				endSum = startSum + savings + revenue
			}

			newCalculation.Revenue = revenue
			newCalculation.EndSum = endSum
			bankDeposit.Calculations = append(bankDeposit.Calculations, newCalculation)

			startMonth = endMonth
			startSum = endSum

			bankDeposit.TotalRevenue += revenue
		}
		bankDeposit.EndSum = endSum

		if startMonth != investmentPeriod {
			diff := investmentPeriod - startMonth
			bankDeposit.EndSum += float64(diff * monthlyPayment)
		}
	}
}

func SortByEndSum(bankDeposits *[]BankDeposit) {
	sort.Slice(*bankDeposits, func(i, j int) bool {
		return (*bankDeposits)[i].EndSum > (*bankDeposits)[j].EndSum
	})
}
