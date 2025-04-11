package entity

type UserData struct {
	LoanId       int    `json:"loan_id"`
	SelfEmployed string `json:"self_employed"`
	IncomeAnnum  int64  `json:"income_annum"`
	LoanAmount   int64  `json:"loan_amount"`
	CibilScore   int    `json:"cibil_score"`
	LoanStatus   string `json:"loan_status"`
}

/* promt
Далее необходимо добавить блок с формой для заполнения данных
Форма содержит следующие поля
Работает ли человек (да/нет)
годовой заработок
сумма кредита
кредитный рейтинг

Далее кнопка получить предикат

по нажатию на кнопку формируется post запрос js к адресу http://localhost:8000/predict в теле запроса должен быть формат json с данными из заполненной формы в виде например
 {
   "loan_id": 6,
   "self_employed": "Yes",
   "income_annum": 4800000,
   "loan_amount": 13500000,
   "cibil_score": 319,
   "loan_status": "Rejected"
  }
в результате строится график снизу
*/
