package utils

import (
	"context"
	"log"
	"wms/config"
	"wms/models"

	"github.com/omniful/go_commons/db/sql/postgres"
	"github.com/omniful/go_commons/redis"
	"gorm.io/gorm"
)

var RedisClient *redis.Client
var PostgresDB *gorm.DB
var Ctx context.Context

func InitializePostgresDb(master postgres.DBConfig, slaves *[]postgres.DBConfig) {
	dbCluster := postgres.InitializeDBInstance(master, slaves)
	Ctx = context.Background()
	PostgresDB = dbCluster.GetMasterDB(Ctx)
	if PostgresDB.Error != nil {
		panic(PostgresDB.Error)
	}
	err := PostgresDB.AutoMigrate(&models.Tenant{}, &models.Seller{}, &models.Hub{}, &models.Sku{}, &models.Inventory{})
	if err != nil {
		panic(err)
	}
	log.Println("Database Connected successfully!")
}

func InitializeRedis() {
	RedisClient = redis.NewClient(config.RedisConfig)
	if err := RedisClient.Ping(Ctx).Err(); err != nil {
		panic(err)
	}
	log.Println("Redis client connected successsfully!")
}
