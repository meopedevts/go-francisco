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

type Customer struct {
	Balance       int       `json:"total"`
	DateStatement time.Time `json:"data_extrato"`
	Limit         int       `json:"limite"`
}

type Transactions struct {
	Amount      int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	CreatedAt   time.Time `json:"realizada_em"`
}

type Statement struct {
	Balance          Customer       `json:"saldo"`
	LastTransactions []Transactions `json:"ultimas_transacoes"`
}

func GetStatement(routing *gin.Engine, db *pgxpool.Pool) {
	routing.GET("/clientes/:id/extrato", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusUnprocessableEntity, "Parametro inválido.")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		var customer Customer
		rowCustomer := db.QueryRow(ctx, "SELECT c.balance, c.limit FROM customers c WHERE c.id = $1", id)
		if err := rowCustomer.Scan(&customer.Balance, &customer.Limit); err != nil {
			if err != nil {
				c.String(http.StatusNotFound, "Usuário não encontrado")
				return
			}
		}
		customer.DateStatement = time.Now()

		rowTransactions, err := db.Query(ctx, "SELECT amount, type, description, created_at FROM transactions WHERE customer_id = $1 ORDER BY created_at DESC LIMIT 10", id)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		defer rowTransactions.Close()

		var lastTransactions []Transactions
		for rowTransactions.Next() {
			var transaction Transactions
			if err := rowTransactions.Scan(&transaction.Amount, &transaction.Type, &transaction.Description, &transaction.CreatedAt); err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			lastTransactions = append(lastTransactions, transaction)
		}
		if err := rowTransactions.Err(); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		statement := Statement{
			customer,
			lastTransactions,
		}
		statementResponse, err := sonic.Marshal(&statement)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.Data(http.StatusOK, "application/json", statementResponse)
	})
}
