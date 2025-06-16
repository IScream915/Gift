package main

import (
	"context"
	"fmt"
	"gift/repo/models"
	"gift/util/mysqlutil"
	"gift/util/redisutil"
	"log"
)

func main() {
	// 根context
	rootCtx := context.Background()

	// 连接Mysql
	db, err := mysqlutil.NewMysqlClient(rootCtx)
	if err != nil {
		log.Fatalln("mysql连接失败", err)
	}
	fmt.Printf("mysqlDb = %#v\n", db)

	// 自动创建表
	if err = db.AutoMigrate(models.Inventory{}, models.Order{}); err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}

	// 连接Redis
	rds, err := redisutil.NewRedisClient(rootCtx)
	if err != nil {
		log.Fatalln("redis连接失败", err)
	}

	fmt.Println(rds)

	// repo层初始化

	// service层初始化

	// api层初始化

}
