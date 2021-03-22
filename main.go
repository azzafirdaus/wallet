package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// init database
	initDB()

	router := httprouter.New()

	// Routes from path to handler function.
	router.POST("/api/v1/init", HandleInitSession)
	router.POST("/api/v1/wallet", Middleware(HandleEnableWallet))
	router.GET("/api/v1/wallet", Middleware(HandleViewBalance))
	router.POST("/api/v1/wallet/deposits", Middleware(HandleDeposits))
	router.POST("/api/v1/wallet/withdrawals", Middleware(HandleWithdrawal))
	router.PATCH("/api/v1/wallet", Middleware(HandleDisableWallet))

	log.Println("starting wallet service at port 8000")

	// Bind to a port and pass router
	log.Fatal(http.ListenAndServe(":8000", router))
}
