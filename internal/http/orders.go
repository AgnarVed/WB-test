package http

import (
	"database/sql"
	"errors"
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

func (h *Handler) createOrder(ctx *fiber.Ctx) error {
	var order models.Order
	err := ctx.BodyParser(order)
	if err != nil {
		h.Response(ctx, http.StatusBadRequest, nil, err)
	}
	err = h.services.Order.CreateOrder(ctx.UserContext(), &order)
	if err != nil {
		return h.Response(ctx, http.StatusInternalServerError, nil, err)
	}
	return h.Response(ctx, http.StatusCreated, nil, nil)
}
