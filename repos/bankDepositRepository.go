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
	Id                  int
	BankName            string
	DepositName         string
	Rate                float64
	NumberOfMonths      int
	Capitalization      bool
	Expandable          bool
	ExpandableLimitFlag bool
	ExpandableLimitKoef float64
	RateLevelUpFlag     bool
	RateLevel2Flag      bool
	RateLevel2Sum       float64
	RateLevel2          float64
	RateLevel3Flag      bool
	RateLevel3Sum       float64
	RateLevel3          float64
	Calculations        []Calculation
	EndSum              float64
	TotalRevenue        float64
}

func GetBankDeposits() *[]BankDeposit {
	rows := database.ExecQuery("SELECT * FROM bank_deposits;")
	defer rows.Close()
	var bankDeposits []BankDeposit
	for rows.Next() {
		var bankDeposit BankDeposit
		err := rows.Scan(&bankDeposit.Id, &bankDeposit.BankName, &bankDeposit.DepositName, &bankDeposit.Rate,
			&bankDeposit.NumberOfMonths, &bankDeposit.Capitalization, &bankDeposit.Expandable,
			&bankDeposit.ExpandableLimitFlag, &bankDeposit.ExpandableLimitKoef, &bankDeposit.RateLevelUpFlag,
			&bankDeposit.RateLevel2Flag, &bankDeposit.RateLevel2Sum, &bankDeposit.RateLevel2,
			&bankDeposit.RateLevel3Flag, &bankDeposit.RateLevel3Sum, &bankDeposit.RateLevel3)
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
			rate := bankDeposit.Rate

			newCalculation := Calculation{
				StartMonth: startMonth,
				EndMonth:   endMonth,
				StartSum:   startSum,
				Savings:    savings,
			}

			if bankDeposit.Expandable {
				preExpandSum := startSum
				startSum += savings
				savings = 0
				if bankDeposit.ExpandableLimitFlag && startSum > preExpandSum*bankDeposit.ExpandableLimitKoef {
					savings = startSum - preExpandSum*bankDeposit.ExpandableLimitKoef
					startSum = preExpandSum * bankDeposit.ExpandableLimitKoef
				}
			}

			if bankDeposit.RateLevelUpFlag {
				if bankDeposit.RateLevel3Flag && startSum >= bankDeposit.RateLevel3Sum {
					rate = bankDeposit.RateLevel3
				} else if bankDeposit.RateLevel2Flag && startSum >= bankDeposit.RateLevel2Sum {
					rate = bankDeposit.RateLevel2
				}
			}

			if bankDeposit.Capitalization {
				revenue = startSum*math.Pow(1+rate/100/12, float64(bankDeposit.NumberOfMonths)) -
					startSum
			} else {
				revenue = (startSum * rate * (depositLength / 366)) / 100
			}
			endSum = startSum + savings + revenue

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
