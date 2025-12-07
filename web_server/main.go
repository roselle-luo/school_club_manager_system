package main

import (
	"log"
	"os"
	"path/filepath"
	"web_server/api"
	"web_server/config"
	"web_server/db"
	"web_server/db/migrate"
	"web_server/docs"
	"web_server/internal/store"
	"web_server/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title 高校社团信息管理系统 API
// @version 1.0
// @description 基于 Vue3+UniApp+Gin+Gorm 的社团管理系统后端接口
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @BasePath /api/v1

func main() {
	cfg := config.Default()
	docs.SwaggerInfo.Title = "高校社团信息管理系统 API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"
	d, err := db.Open(cfg.DB)
	if err != nil {
		logger.Error("db open error:", err)
		log.Fatal(err)
	}
	store.SetDB(d)
	if err := migrate.AutoMigrate(d); err != nil {
		logger.Error("auto migrate error:", err)
		log.Fatal(err)
	}
	_ = migrate.Seed(d)

	r := gin.Default()
	pubPath := filepath.Join(cfg.Server.PublicDir)
	upPath := filepath.Join(pubPath, cfg.Server.UploadDir)
	_ = os.MkdirAll(upPath, 0755)
	r.Static("/static", pubPath)
	api.Register(r)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(cfg.Server.Addr); err != nil {
		logger.Error("server run error:", err)
		log.Fatal(err)
	}
}
