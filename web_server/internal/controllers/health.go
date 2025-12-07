package controllers

import (
	"github.com/gin-gonic/gin"
	"web_server/pkg/response"
)

// @Summary 健康检查
// @Tags 公共
// @Success 200 {object} response.Body
// @Router /health [get]
func Health(c *gin.Context) {
	c.JSON(200, response.Success(map[string]string{"status": "ok"}))
}
