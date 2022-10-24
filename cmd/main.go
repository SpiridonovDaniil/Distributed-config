package main

import (
	router "github.com/SpiridonovDaniil/Distributed-config/internal/app/http"
	"github.com/SpiridonovDaniil/Distributed-config/internal/app/service"
	"github.com/SpiridonovDaniil/Distributed-config/internal/repository/postgres"
)

func main() {
	db := postgres.New("db:5432", "user", "test", "config")
	service := service.New(db)
	r := router.NewServer(service)
	err := r.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
