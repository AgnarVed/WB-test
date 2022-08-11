package http

import (
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
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
