package models

type Response struct {
	ErrorText string      `json:"errorText"`
	HasError  bool        `json:"hasError"`
	Resp      interface{} `json:"resp"`
}
