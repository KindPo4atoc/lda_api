package repository

import (
	"github.com/lda_api/internal/app/entity"
)

// структура для взаимодействия с бд
type DataForLearnRepository struct {
	store *DataBase
}

func (r *DataForLearnRepository) SelectAllLearnData() (entity.ContextData, error) {
	var users entity.ContextData
	rows, err := r.store.db.Query(
		"SELECT loan_id, loan_term, income_annum, loan_amount, cibil_score, loan_status " +
			"FROM learndata where type_data = 0;",
	)

	if err != nil {
		return entity.ContextData{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var dataRow entity.UserData

		err := rows.Scan(
			&dataRow.LoanId,
			&dataRow.LoanTerm,
			&dataRow.IncomeAnnum,
			&dataRow.LoanAmount,
			&dataRow.CibilScore,
			&dataRow.LoanStatus,
		)
		if err != nil {
			return entity.ContextData{}, err
		}
		users.Data = append(users.Data, dataRow)
	}

	return users, nil
}

func (r *DataForLearnRepository) SelectingDataByClass(classData string) (entity.ContextData, error) {
	var dataByClass entity.ContextData

	rows, err := r.store.db.Query(
		"SELECT loan_id, self_employed, income_annum, loan_amount, cibil_score, loan_status "+
			"FROM learndata where loan_status = $1;",
		classData,
	)

	if err != nil {
		return entity.ContextData{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var dataRow entity.UserData
		err := rows.Scan(
			&dataRow.LoanId,
			&dataRow.LoanTerm,
			&dataRow.IncomeAnnum,
			&dataRow.LoanAmount,
			&dataRow.CibilScore,
			&dataRow.LoanStatus,
		)
		if err != nil {
			return entity.ContextData{}, err
		}
		dataByClass.Data = append(dataByClass.Data, dataRow)
	}

	return dataByClass, nil
}
func (r *DataForLearnRepository) SelectTestData() (entity.ContextData, error) {
	var users entity.ContextData
	rows, err := r.store.db.Query(
		"SELECT loan_id, loan_term, income_annum, loan_amount, cibil_score, loan_status " +
			"FROM learndata where type_data = 1;",
	)

	if err != nil {
		return entity.ContextData{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var dataRow entity.UserData

		err := rows.Scan(
			&dataRow.LoanId,
			&dataRow.LoanTerm,
			&dataRow.IncomeAnnum,
			&dataRow.LoanAmount,
			&dataRow.CibilScore,
			&dataRow.LoanStatus,
		)
		if err != nil {
			return entity.ContextData{}, err
		}
		users.Data = append(users.Data, dataRow)
	}

	return users, nil
}
