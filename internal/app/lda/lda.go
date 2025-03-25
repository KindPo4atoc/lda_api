package lda

import (
	"github.com/lda_api/internal/app/entity"
)

type LDA struct {
	Alpha float64   `json:"alpha"`
	Beta  float64   `json:"beta"`
	X     []float64 `json:"X"`
	Y     []float64 `json:"Y"`
}

func (lda *LDA) FitModel(usr []entity.UserData) {

}

func (lda *LDA) Mean() float64 {

	return 0.0
}
