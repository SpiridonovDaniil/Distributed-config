package main

import (
	router "github.com/SpiridonovDaniil/Distributed-config/internal/app/http"
	"github.com/SpiridonovDaniil/Distributed-config/internal/app/service"
	"github.com/SpiridonovDaniil/Distributed-config/internal/config"
	"github.com/SpiridonovDaniil/Distributed-config/internal/repository/postgres"
)

func main() {
	cfg := config.Read()

	db := postgres.New(cfg.Postgres)
	service := service.New(db)
	r := router.NewServer(service)
	err := r.Listen(":" + cfg.Service.Port)
	if err != nil {
		panic(err)
	}
}
