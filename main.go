package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
)

var db *gorm.DB

type Todo struct {
	gorm.Model
	Title   string `form:"title" json:"title"`
	Content string `form:"content" json:"content"`
	Status  bool   `form:"status" json:"status"`
}

func main() {

	err := initDB()
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Todo{})

	//注册路由
	r := gin.Default()
	g := r.Group("/api/v1")
	{
		g.POST("/create", createTodo)
		g.PUT("/update", updateTodo)
		g.GET("/get", getTodo)
		g.DELETE("/delete/:id", deleteTodo)

	}

	r.Run(":8089")

}

func initDB() (err error) {
	dsn := "root:root@tcp(localhost:3306)/Todo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}

func createTodo(c *gin.Context) {
	//1获取请求参数
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "无效的参数",
		})
		return
	}
	//2处理业务逻辑
	if err := db.Create(&todo).Error; err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "创建失败",
		})
		return
	}
	//3返回响应
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": todo,
	})
}

func updateTodo(c *gin.Context) {
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "无效的参数",
		})
		return
	}
	//业务逻辑
	if err := db.First(&Todo{}, todo.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "记录不存在",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "服务端错误",
			})
		}
		return
	}
	if err := db.Model(&todo).Update("status", todo.Status).Error; err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "更新失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "更新成功",
	})
}

func getTodo(c *gin.Context) {
	var todos []Todo
	if err := db.Find(&todos).Error; err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "服务端错误",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": todos,
	})
}

func deleteTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "无效的参数",
		})
		return
	}
	if err := db.First(&Todo{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "记录不存在",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "服务端错误",
			})
			return
		}
	}
	if err := db.Delete(&Todo{}, id).Error; err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "服务端错误",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
	})

}
