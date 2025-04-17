package lda

// структура для хранения и отправки предикта
type PredictModel struct {
	Predict        int         `json:"predict_model"`
	Distance       [][]float64 `json:"score_discrimination"`
	ConvertionData [][]float64 `json:"data"`
}

func NewPredict(s int, d [][]float64, data [][]float64) *PredictModel {
	return &PredictModel{
		Predict:        s,
		Distance:       d,
		ConvertionData: data,
	}
}
