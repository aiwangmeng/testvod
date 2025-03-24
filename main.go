package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 数据表结构
type Movie struct {
	ID         uint      `gorm:"primaryKey"`
	Name       string    `gorm:"type:varchar(255)"`
	TypeIDs    string    `gorm:"type:varchar(255)"`
	ClassNames string    `gorm:"type:varchar(255)"`
	Tags       string    `gorm:"type:varchar(255)"`
	Level      int       `gorm:"type:tinyint"`
	Area       string    `gorm:"type:varchar(100)"`
	Lang       string    `gorm:"type:varchar(100)"`
	Year       int       `gorm:"type:year(4)"`
	State      string    `gorm:"type:varchar(50)"`
	Version    string    `gorm:"type:varchar(50)"`
	Weekday    int       `gorm:"type:tinyint"`
	RelIDs     string    `gorm:"type:varchar(255)"`
	TimeAdd    time.Time `gorm:"type:datetime"`
	TimeHits   time.Time `gorm:"type:datetime"`
	Time       time.Time `gorm:"type:datetime"`
	Hits       int       `gorm:"type:int unsigned"`
	HitsDay    int       `gorm:"type:int unsigned"`
	HitsWeek   int       `gorm:"type:int unsigned"`
	HitsMonth  int       `gorm:"type:int unsigned"`
	Up         int       `gorm:"type:int unsigned"`
	Down       int       `gorm:"type:int unsigned"`
	IsEnd      bool      `gorm:"type:tinyint(1)"`
	Plot       bool      `gorm:"type:tinyint(1)"`
}

// QueryParams 请求参数结构
type QueryParams struct {
	Order     string `form:"order"`
	By        string `form:"by"`
	Start     int    `form:"start"`
	Num       int    `form:"num"`
	IDs       string `form:"ids"`
	Not       string `form:"not"`
	Type      string `form:"type"`
	Class     string `form:"class"`
	Tag       string `form:"tag"`
	Level     string `form:"level"`
	Area      string `form:"area"`
	Lang      string `form:"lang"`
	Year      string `form:"year"`
	State     string `form:"state"`
	Version   string `form:"version"`
	Weekday   string `form:"weekday"`
	Rel       string `form:"rel"`
	TimeAdd   string `form:"timeadd"`
	TimeHits  string `form:"timehits"`
	Time      string `form:"time"`
	HitsMonth string `form:"hitsmonth"`
	HitsWeek  string `form:"hitsweek"`
	HitsDay   string `form:"hitsday"`
	Hits      string `form:"hits"`
	IsEnd     string `form:"isend"`
	Plot      string `form:"plot"`
}

func main() {
	// 初始化Gin
	r := gin.Default()

	// 初始化数据库连接
	dsn := "root:@tcp(127.0.0.1:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(errors.New(err.Error()))
		return
	}

	// 注册路由
	r.GET("/select", func(c *gin.Context) {
		var params QueryParams
		err := c.ShouldBindQuery(&params)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := db.Model(&Movie{})

		if params.IDs != "" {
			query.Where("id IN (?)", strings.Split(params.IDs, ","))
		}
		if params.Not != "" {
			query.Where("id NOT IN (?)", strings.Split(params.Not, ","))
		}

		if params.Type != "all" && params.Type != "" {
			query.Where("type_ids like ?", fmt.Sprintf("%s%", params.Type))
		}

		if params.Tag != "" {
			tags := strings.Split(params.Tag, ",")
			for _, tag := range tags {
				query.Where("FIND_IN_SET(?, tags)", tag)
			}
		}

		if params.Level != "" {
			query.Where("level IN (?)", strings.Split(params.Level, ","))
		}

		if params.IsEnd != "" {
			query.Where("is_end = ?", params.IsEnd)
		}

		if params.By != "" {
			query.Order(clause.OrderByColumn{
				Column: clause.Column{Name: params.By},
				Desc:   params.Order == "desc",
			})
		}
		//分页
		query.Offset(params.Start).Limit(params.Num)

		// 执行查询
		var results []Movie
		err = query.Find(&results).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data":    results,
			"total":   len(results),
		})
	})

	// 启动服务
	r.Run(":8082")
}
