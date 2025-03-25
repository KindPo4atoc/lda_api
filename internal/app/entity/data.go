package entity

type UserData struct {
	LoanId       int    `json:"loan_id"`
	SelfEmployed string `json:"self_employed"`
	IncomeAnnum  int64  `json:"income_annum"`
	LoanAmount   int64  `json:"loan_amount"`
	CibilScore   int    `json:"cibil_score"`
	LoanStatus   string `json:"loan_status"`
}
