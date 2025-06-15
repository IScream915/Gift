package main

import (
	"gift/repo/models"
	"gift/util/mysqlutil"
	"log"
)

func main() {
	//fmt.Println("Hello World")
	db, err := mysqlutil.NewMysqlClient()
	if err != nil {
		log.Fatalln("数据库连接失败", err)
	}
	//fmt.Printf("mysqlDb = %#v\n", db)
	// 自动创建表
	if err = db.AutoMigrate(models.Inventory{}); err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}

	// repo层初始化

	// service层初始化

	// api层初始化

}
