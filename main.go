package main

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sergazyyev/wallet/app"
	"log"
)

func main() {
	log.Println(app.Start())
}
