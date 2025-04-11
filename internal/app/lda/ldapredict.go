package lda

import "github.com/lda_api/internal/app/entity"

type PredictModel struct {
	Predict             string                `json:"predict_model"`
	ScoreDiscrimitation float64               `json:"score_discrimination"`
	ConvertionData      entity.ConversionData `json:"data"`
}

func NewPredict(s string, d float64, data entity.ConversionData) *PredictModel {
	return &PredictModel{
		Predict:             s,
		ScoreDiscrimitation: d,
		ConvertionData:      data,
	}
}
