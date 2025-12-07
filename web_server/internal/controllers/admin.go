package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"web_server/db/models"
	"web_server/internal/store"
	"web_server/pkg/response"
)

type UpdateRoleReq struct {
	Role string `json:"role" binding:"required"`
}

// @Summary 更新成员社团内角色
// @Tags 管理员
// @Accept json
// @Produce json
// @Param id path int true "成员关系ID"
// @Param payload body UpdateRoleReq true "角色值"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /admin/memberships/{id}/role [post]
func UpdateMembershipRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	var req UpdateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	switch req.Role {
	case "member", "leader", "advisor":
	default:
		c.JSON(http.StatusBadRequest, response.Error(400, "非法角色"))
		return
	}
	var m models.Membership
	if err := store.DB().Where("id = ?", id).First(&m).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "不存在"))
		return
	}
	m.Role = req.Role
	if err := store.DB().Save(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "更新失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(m))
}
