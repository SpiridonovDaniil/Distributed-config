package main

import (
	"github.com/SpiridonovDaniil/Distributed-config/internal/app/http"
	"github.com/SpiridonovDaniil/Distributed-config/internal/app/service"
	"github.com/SpiridonovDaniil/Distributed-config/internal/repository/postgres"
)

func main() {
	db := postgres.New("", "", "", "")
	service := service.New(db)
	http.NewServer(service)
}
