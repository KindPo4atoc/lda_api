package entity

type ContextData struct {
	Data []UserData `json:"users_data"`
}

func New(d []UserData) *ContextData {
	return &ContextData{
		Data: d,
	}
}
