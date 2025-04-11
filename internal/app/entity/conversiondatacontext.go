package entity

type ConvertDataContext struct {
	Data []ConversionData `json:"convesion_data"`
}

func NewConvertContext(d []ConversionData) *ConvertDataContext {
	return &ConvertDataContext{
		Data: d,
	}
}
