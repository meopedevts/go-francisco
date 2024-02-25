package routes

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RequestTransaction struct {
	Value       int    `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

type TransactionResponse struct {
	Limit   int `json:"limite"`
	Balance int `json:"saldo"`
}

func PostTransactions(routing *gin.Engine, db *pgxpool.Pool) {
	routing.POST("/clientes/:id/transacoes", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusUnprocessableEntity, "Parametro inválido.")
			return
		}

		var requestTransaction RequestTransaction
		if err := c.BindJSON(&requestTransaction); err != nil {
			c.String(http.StatusUnprocessableEntity, "Requisição inválida.")
			return
		}

		if requestTransaction.Type != "c" && requestTransaction.Type != "d" {
			c.String(http.StatusUnprocessableEntity, "Tipo inválido.")
			return
		}

		if len(requestTransaction.Description) > 10 || len(requestTransaction.Description) < 1 {
			c.String(http.StatusUnprocessableEntity, "Descrição inválida.")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		var customer Customer
		rowCustomer := db.QueryRow(ctx, "SELECT c.balance, c.limit  FROM customers c WHERE c.id = $1", id)
		if err := rowCustomer.Scan(&customer.Balance, &customer.Limit); err != nil {
			if err != nil {
				c.String(http.StatusNotFound, "Usuário não encontrado")
				return
			}
		}

		var newBalance int
		if requestTransaction.Type == "c" {
			newBalance = customer.Balance + requestTransaction.Value
		}
		if requestTransaction.Type == "d" {
			newBalance = customer.Balance - requestTransaction.Value
		}
		if newBalance < (customer.Limit * -1) {
			c.String(http.StatusUnprocessableEntity, "Limite insuficiente.")
			return
		}

		_, err = db.Exec(ctx, "INSERT INTO transactions (customer_id, amount, type, description, created_at) VALUES ($1, $2, $3, $4, $5)",
			id, requestTransaction.Value, requestTransaction.Type, requestTransaction.Description, time.Now(),
		)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		_, err = db.Exec(ctx, "UPDATE customers SET balance = $1 WHERE id = $2", newBalance, id)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		transactionResponse := TransactionResponse{
			customer.Limit,
			newBalance,
		}
		parsedResponse, err := sonic.Marshal(&transactionResponse)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.Data(http.StatusOK, "application/json", parsedResponse)
	})
}
