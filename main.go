package main

import (
	"github.com/meopedevts/go-francisco/db"
	"github.com/meopedevts/go-francisco/routing"
)

func main() {
	conn := db.StartDb()

	routing.Open(conn)
}
