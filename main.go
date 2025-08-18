package main

import (
	"log"
	"net/http"

	"tracker/config"
	"tracker/database"
	"tracker/handler"
	"tracker/repository"
	"tracker/routes"
	"tracker/service"
)

func main() {
	// 1) env
	config.LoadEnv() // or handle the error if you changed LoadEnv to return error

	// 2) db
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	// 3) repos (with DB fields added)
	userRepo := &repository.UserRepo{DB: db}
	txRepo   := &repository.TransactionRepo{DB: db}
	budRepo  := &repository.BudgetRepo{DB: db}

	// 4) services
	userSvc := &service.UserService{Repo: userRepo}
	txSvc   := &service.TransactionService{Repo: txRepo}
	budSvc  := &service.BudgetService{Repo: budRepo}

	// 5) handlers
	userH := &handler.UserHandler{Service: userSvc}
	txH   := &handler.TransactionHandler{Service: txSvc}
	budH  := &handler.BudgetHandler{Service: budSvc}

	// 6) router
	r := routes.SetupRouter(userH, txH, budH)

	log.Println("listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}


