package controllers

import (
	"net/http"
	"strconv"
	"web_server/db/models"
	"web_server/internal/authz"
	"web_server/internal/store"
	"web_server/pkg/response"

	"github.com/gin-gonic/gin"
)

// @Summary 获取社团负责人列表
// @Tags 负责人
// @Produce json
// @Param clubId path int true "社团ID"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/clubs/{clubId}/users [get]
func GetClubLeaders(c *gin.Context) {
	clubIDStr := c.Param("clubId")
	clubID, err := strconv.Atoi(clubIDStr)
	if err != nil || clubID <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	cu, _ := c.Get("currentUser")
	user := cu.(*models.User)
	if !(authz.IsAdmin(user) || authz.IsClubLeader(user.ID, uint(clubID))) {
		c.JSON(http.StatusForbidden, response.Error(403, "无权限"))
		return
	}
	var items []models.Membership
	if err := store.DB().Where("club_id = ? AND role IN ?", clubID, []string{"leader", "advisor"}).Preload("User").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
		return
	}
	users := make([]models.User, 0, len(items))
	for _, m := range items {
		users = append(users, m.User)
	}
	c.JSON(http.StatusOK, response.Success(users))
}

// @Summary 获取用户负责的社团列表
// @Tags 负责人
// @Produce json
// @Param userId path int true "用户ID"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/users/{userId}/clubs [get]
func GetUserLeaderClubs(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	cu, _ := c.Get("currentUser")
	user := cu.(*models.User)
	if !(authz.IsAdmin(user) || user.ID == uint(userID)) {
		c.JSON(http.StatusForbidden, response.Error(403, "无权限"))
		return
	}
	var items []models.Membership
	if err := store.DB().Where("user_id = ? AND role IN ?", userID, []string{"leader", "advisor"}).Preload("Club").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
		return
	}
	clubs := make([]models.Club, 0, len(items))
	for _, m := range items {
		clubs = append(clubs, m.Club)
	}
	c.JSON(http.StatusOK, response.Success(clubs))
}

type setRoleReq struct {
	Role string `json:"role" binding:"required"`
}

// @Summary 负责人设定成员社团内角色
// @Tags 负责人
// @Accept json
// @Produce json
// @Param clubId path int true "社团ID"
// @Param userId path int true "用户ID"
// @Param payload body setRoleReq true "角色值"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/clubs/{clubId}/members/{userId}/role [post]
func SetMemberRoleByLeader(c *gin.Context) {
	clubIDStr := c.Param("clubId")
	userIDStr := c.Param("userId")
	clubID, err1 := strconv.Atoi(clubIDStr)
	userID, err2 := strconv.Atoi(userIDStr)
	if err1 != nil || err2 != nil || clubID <= 0 || userID <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	cu, _ := c.Get("currentUser")
	user := cu.(*models.User)
	if !(authz.IsAdmin(user) || authz.IsClubLeader(user.ID, uint(clubID))) {
		c.JSON(http.StatusForbidden, response.Error(403, "无权限"))
		return
	}
	var req setRoleReq
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
	if err := store.DB().Where("user_id = ? AND club_id = ?", userID, clubID).First(&m).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "成员不存在该社团"))
		return
	}
	m.Role = req.Role
	if err := store.DB().Save(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "更新失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(m))
}
