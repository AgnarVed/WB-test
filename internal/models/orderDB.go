package models

import "encoding/json"

type Order struct {
	ID       int64           `json:"ID"`
	OrderUID string          `json:"order_uid"`
	Data     json.RawMessage `json:"data"`
}
