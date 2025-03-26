package repository

import (
	"github.com/lda_api/internal/app/entity"
)

type DataForLearnRepository struct {
	store *DataBase
}

func (r *DataForLearnRepository) SelectAllData() ([]entity.UserData, error) {
	var users []entity.UserData
	rows, err := r.store.db.Query("SELECT * FROM learndata limit 30;")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dataRow entity.UserData

		err := rows.Scan(
			&dataRow.LoanId,
			&dataRow.SelfEmployed,
			&dataRow.IncomeAnnum,
			&dataRow.LoanAmount,
			&dataRow.CibilScore,
			&dataRow.LoanStatus,
		)
		if err != nil {
			return nil, err
		}
		if dataRow.LoanId == 29 {
			dataRow.LoanId = 34
			dataRow.SelfEmployed = " Yes"
			dataRow.IncomeAnnum = 8400000
			dataRow.LoanAmount = 22000000
			dataRow.CibilScore = 830
			dataRow.LoanStatus = " Approved"
		}
		users = append(users, dataRow)
	}

	return users, nil
}

func (r *DataForLearnRepository) SelectingDataByClass(classData string) ([]entity.UserData, error) {
	var dataByClass []entity.UserData

	rows, err := r.store.db.Query("SELECT * FROM learndata where loan_status = $1 limit 15;", classData)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dataRow entity.UserData
		err := rows.Scan(
			&dataRow.LoanId,
			&dataRow.SelfEmployed,
			&dataRow.IncomeAnnum,
			&dataRow.LoanAmount,
			&dataRow.CibilScore,
			&dataRow.LoanStatus,
		)
		if err != nil {
			return nil, err
		}
		dataByClass = append(dataByClass, dataRow)
	}

	return dataByClass, nil
}
