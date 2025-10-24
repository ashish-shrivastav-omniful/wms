package main

import (
	"wms/config"
	"wms/routes"
	"wms/utils"

	"github.com/omniful/go_commons/db/sql/postgres"
	"github.com/omniful/go_commons/http"
)

func main() {
	server := http.InitializeServer(":8080", 0, 0, 0, false)
	utils.InitializePostgresDb(*config.PostgresConfig, &[]postgres.DBConfig{*config.PostgresConfig})
	utils.InitializeRedis()
	routes.GetRoutes(server)
	if err := server.StartServer("wms"); err != nil {
		panic(err.Error())
	}
}
