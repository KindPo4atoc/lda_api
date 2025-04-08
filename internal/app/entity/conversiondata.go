package entity

type ConversionData struct {
	ImportancecCoefficient float64 `json:"coefficient"`
	Rating                 float64 `json:"rating"`
	Class                  int     `json:"class"`
}
