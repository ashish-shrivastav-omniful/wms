package config

import (
	"github.com/omniful/go_commons/db/sql/postgres"
	"github.com/omniful/go_commons/redis"
)

var RedisConfig = &redis.Config{
	ClusterMode: false,
	Hosts:       []string{"localhost:6379"},
	PoolSize:    1,
}
var PostgresConfig = &postgres.DBConfig{
	Host:     "localhost",
	Port:     "5432",
	Username: "root",
	Password: "ashish",
	Dbname:   "wms",
}
