package models

type Balance struct{
	Symbol string
	Locked string
	Free string
	Total string
}

type UserBalances struct {
	Balances []Balance
}

type UserBalanceData struct{
	UserId string
	Balances UserBalances
}



