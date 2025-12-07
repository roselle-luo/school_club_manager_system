package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web_server/config"
	"web_server/db/models"
	"web_server/internal/store"
	"web_server/pkg/jwt"
	"web_server/pkg/password"
	"web_server/pkg/response"
)

type LoginReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary 登录
// @Tags 公共
// @Accept json
// @Produce json
// @Param payload body LoginReq true "登录参数"
// @Success 200 {object} response.Body
// @Router /public/login [post]
func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	var u models.User
	if err := store.DB().Where("account = ?", req.Account).Preload("Role").First(&u).Error; err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(401, "账号或密码错误"))
		return
	}
	if !password.Compare(u.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, response.Error(401, "账号或密码错误"))
		return
	}
	token, err := jwt.Sign(config.Default().JWT.Secret, u.ID, u.Role.Code, config.Default().JWT.Expires)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "登录失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(map[string]any{"token": token}))
}
