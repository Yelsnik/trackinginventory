package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/Yelsnik/trackinginventory/db/sqlc"
	"github.com/Yelsnik/trackinginventory/token"

	"github.com/gin-gonic/gin"

	"github.com/lib/pq"
)

type createInventoryRequest struct {
	Item     string `json:"item" binding:"required"`
	SerialNo string `json:"serialno" binding:"required"`
	Price    int64  `json:"price" binding:"required,min=1"`
}

func (server *Server) createInventory(ctx *gin.Context) {
	var req createInventoryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateInventoryParams{
		Item:         req.Item,
		SerialNumber: req.SerialNo,
		Price:        req.Price,
		Owner:        authPayload.Owner,
	}

	fmt.Println(authPayload)

	inventory, err := server.db.CreateInventory(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {

			switch pqErr.Error() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}

		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, inventory)
}

type getInventoryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getInventory(ctx *gin.Context) {
	var req getInventoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	inventory, err := server.db.GetInventory(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if inventory.Owner != authPayload.Owner {
		err := errors.New("account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, inventory)
}
