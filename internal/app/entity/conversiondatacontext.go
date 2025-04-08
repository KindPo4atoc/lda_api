package entity

type ConvertDataContext struct {
	Data []ConversionData `json:"convesion__data"`
}

func NewConvertContext(d []ConversionData) *ConvertDataContext {
	return &ConvertDataContext{
		Data: d,
	}
}
