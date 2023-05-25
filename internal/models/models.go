package models

import "encoding/json"

type OrderInput struct {
	OrderUID  string          `json:"order_uid"`
	OrderInfo json.RawMessage `json:"order_info"`
}
