package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meopedevts/go-francisco/routes"
)

func Open(db *pgxpool.Pool) {
	routing := gin.Default()
	routes.GetStatement(routing, db)
	routes.PostTransactions(routing, db)

	routing.Run(":8080")
}
