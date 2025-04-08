package lda

type PredictModel struct {
	Predict             string  `json:"predict_model"`
	ScoreDiscrimitation float64 `json:"score_discrimination"`
}

func NewPredict(s string, d float64) *PredictModel {
	return &PredictModel{
		Predict:             s,
		ScoreDiscrimitation: d,
	}
}
