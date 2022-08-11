package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
	"wbTest/internal/models"
)

func (h *Handler) getOrderByID(ctx *fiber.Ctx) error {
	orderIdStr := ctx.Params("orderID")
	orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
	if err != nil {
		return h.Response(ctx, http.StatusBadRequest, nil, err)
	}

	order, err := h.services.Order.GetOrderByID(ctx.UserContext(), int(orderId))
	if err != nil {
		if err == sql.ErrNoRows {
			return h.Response(ctx, http.StatusOK, nil, errors.New("order not found"))
		}
		return h.Response(ctx, http.StatusInternalServerError, nil, err)
	}
	return h.Response(ctx, http.StatusOK, order, nil)

}

type Msg struct {
	Message string `json:"msg"`
}

// to fix json error while importing to psql
func (h *Handler) showData(ctx *fiber.Ctx) error {
	var data models.Order
	err := ctx.BodyParser(&data)
	if err != nil {
		h.Response(ctx, http.StatusBadRequest, nil, err)
	}
	var dataJson Msg
	err = json.Unmarshal(data.Data, &dataJson)
	if err != nil {
		h.Response(ctx, http.StatusBadRequest, nil, err)
	}

	//newData, err := json.Marshal(&dataJson)
	//if err != nil {
	//	return h.Response(ctx, http.StatusBadRequest, nil, err)
	//}

	fmt.Printf("data json is: %s", dataJson.Message)
	return nil
}

func (h *Handler) createOrder(ctx *fiber.Ctx) error {
	var order models.Order
	err := ctx.BodyParser(&order)
	if err != nil {
		h.Response(ctx, http.StatusBadRequest, nil, err)
	}
	err = h.services.Order.CreateOrder(ctx.UserContext(), &order)
	if err != nil {
		return h.Response(ctx, http.StatusInternalServerError, nil, err)
	}
	return h.Response(ctx, http.StatusCreated, nil, nil)
}
