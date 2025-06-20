package main

import (
	"context"
	"fmt"
	"gift/api"
	"gift/repo"
	"gift/repo/models"
	"gift/services"
	"gift/util/mysqlutil"
	"gift/util/redisutil"
	"github.com/gin-gonic/gin"
	"log"
)

// SetRoute 路由组设置
func SetRoute(router gin.IRouter, inventoryApi api.Inventory) {
	v1 := router.Group("/v1")
	{
		inventoryNotLogin := v1.Group("/inventory")
		{
			inventoryNotLogin.POST("/create", inventoryApi.Create) // 增加新物品
			inventoryNotLogin.POST("/update", inventoryApi.Update) // 更新物品信息
			inventoryNotLogin.POST("/delete", inventoryApi.Delete) // 删除物品
			inventoryNotLogin.GET("/get", inventoryApi.Get)        // 获取物品
			inventoryNotLogin.GET("/load", inventoryApi.Load)      // 数据预热
		}
	}
}

// NewEngine 初始化新engine
func NewEngine(inventoryApi api.Inventory) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	// 启用自定义的引擎
	engine := gin.New()

	// TODO 注册全局中间件
	// gin.Recovery是内置的回复中间件，用于捕获程序中的panic，防止服务崩溃，并返回500错误
	//engine.Use(gin.Recovery(), middleware.RequestID(), middleware.GinLogger())
	engine.Use(gin.Logger(), gin.Recovery())

	// 设置路由
	SetRoute(
		engine,
		inventoryApi,
	)
	return engine
}

func main() {
	// 根context
	rootCtx := context.Background()

	///////////////////////////////////////////// 连接Mysql /////////////////////////////////////////////////
	db, err := mysqlutil.NewMysqlClient(rootCtx)
	if err != nil {
		log.Fatalln("mysql连接失败", err)
	}
	fmt.Println("Mysql初始化成功")
	// 自动创建表
	if err = db.AutoMigrate(models.Inventory{}, models.Order{}); err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}
	////////////////////////////////////////////////////////////////////////////////////////////////////////
	// -> db 全局mysql数据库

	///////////////////////////////////////////// 连接Redis /////////////////////////////////////////////////
	rds, err := redisutil.NewRedisClient(rootCtx)
	if err != nil {
		log.Fatalln("redis连接失败", err)
	}

	fmt.Println("Redis初始化成功", rds)
	////////////////////////////////////////////////////////////////////////////////////////////////////////
	// -> rds全局redis数据库

	/////////////////////////////////////////// 业务逻辑初始化 ///////////////////////////////////////////////
	// repo层初始化
	inventoryRepo := repo.NewInventory(db)
	inventoryRdsRepo := repo.NewInventoryRds(rds)

	// service层初始化
	inventorySvc := services.NewInventory(inventoryRepo, inventoryRdsRepo)

	// api层初始化
	inventoryApi := api.NewInventory(inventorySvc)
	////////////////////////////////////////////////////////////////////////////////////////////////////////

	// 实例化gin引擎
	engine := NewEngine(inventoryApi)
	// 在本地8080端口开启服务
	_ = engine.Run(":8000")
}
