package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	// ID    uint `gorm:"primaryKey,autoIncrement"`
	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	// db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/your_database?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db = db.Debug() // 开启调试模式，打印 SQL 语句

	// 迁移 schema
	// 建表
	db.AutoMigrate(&Product{})

	// Create
	// 插入
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	// 根据整型主键查找
	// 类似 SELECT * FROM products WHERE id = 1 ORDER BY id LIMIT 1;
	db.First(&product, 1)
	// 查找 code 字段值为 D42 的记录
	// 类似 SELECT * FROM products WHERE code = 'D42' ORDER BY id LIMIT 1;
	db.First(&product, "code = ?", "D42")

	// Update - 将 product 的 price 更新为 200
	db.Model(&product).Update("Price", 200)
	// Update - 更新多个字段
	// 这句话会更新两个字段
	// 类似 UPDATE products SET price = 200, code = 'F42' WHERE id = 1;
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	db.Delete(&product, 1)
}
