package repository

import (
	"github.com/lda_api/internal/app/entity"
)

type DataForLearnRepository struct {
	store *DataBase
}

func (r *DataForLearnRepository) SelectAllData() ([]entity.UserData, error) {
	var users []entity.UserData
	rows, err := r.store.db.Query("SELECT * FROM learndata;")

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
		users = append(users, dataRow)
	}

	return users, nil
}
