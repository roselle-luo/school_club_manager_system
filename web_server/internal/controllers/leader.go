package controllers

import (
	"net/http"
	"strconv"
	"web_server/db/models"
	"web_server/internal/authz"
	"web_server/internal/store"
	"web_server/pkg/pagination"
	"web_server/pkg/response"

	"github.com/gin-gonic/gin"
)

// @Summary 获取社团负责人列表
// @Tags 负责人
// @Produce json
// @Param clubId path int true "社团ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
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
	q := store.DB().Model(&models.Membership{}).Where("club_id = ? AND role IN ?", clubID, []string{"leader", "advisor"}).Preload("User").Order("id DESC")
	pg := pagination.Get(c)
	info, err := pagination.Do(q, pg, &items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
		return
	}
	users := make([]models.User, 0, len(items))
	for _, m := range items {
		users = append(users, m.User)
	}
	c.JSON(http.StatusOK, response.Success(map[string]any{"list": users, "pagination": info}))
}

// @Summary 获取用户负责的社团列表
// @Tags 负责人
// @Produce json
// @Param userId path int true "用户ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param keyword query string false "关键词：社团名称"
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

	keyword := c.Query("keyword")

	type ClubWithRole struct {
		models.Club
		Role string `json:"role"`
	}

	// Special case: Admin viewing their own list sees ALL clubs
	if authz.IsAdmin(user) && user.ID == uint(userID) {
		var clubs []models.Club
		q := store.DB().Model(&models.Club{})
		if keyword != "" {
			q = q.Where("name LIKE ?", "%"+keyword+"%")
		}
		q = q.Preload("Category").Order("id DESC") // Preload Category for better display if needed

		pg := pagination.Get(c)
		info, err := pagination.Do(q, pg, &clubs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
			return
		}

		result := make([]ClubWithRole, 0, len(clubs))
		for _, cl := range clubs {
			result = append(result, ClubWithRole{Club: cl, Role: "admin"})
		}
		c.JSON(http.StatusOK, response.Success(map[string]any{"list": result, "pagination": info}))
		return
	}

	var items []models.Membership
	q := store.DB().Model(&models.Membership{}).Where("user_id = ? AND role IN ?", userID, []string{"leader", "advisor"})

	if keyword != "" {
		q = q.Joins("JOIN clubs ON clubs.id = memberships.club_id").Where("clubs.name LIKE ?", "%"+keyword+"%")
	}

	q = q.Preload("Club").Preload("Club.Category").Order("memberships.id DESC")

	pg := pagination.Get(c)
	info, err := pagination.Do(q, pg, &items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
		return
	}
	clubs := make([]ClubWithRole, 0, len(items))
	for _, m := range items {
		clubs = append(clubs, ClubWithRole{Club: m.Club, Role: m.Role})
	}
	c.JSON(http.StatusOK, response.Success(map[string]any{"list": clubs, "pagination": info}))
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

	// 权限检查：除了管理员，普通负责人（社长/管理员）需要检查级别
	if !authz.IsAdmin(user) {
		var callerM models.Membership
		if err := store.DB().Where("user_id = ? AND club_id = ?", user.ID, clubID).First(&callerM).Error; err != nil {
			c.JSON(http.StatusForbidden, response.Error(403, "无权限"))
			return
		}

		getRoleLevel := func(role string) int {
			switch role {
			case "leader":
				return 3
			case "advisor":
				return 2
			default:
				return 1
			}
		}

		myLevel := getRoleLevel(callerM.Role)
		targetLevel := getRoleLevel(m.Role)
		newLevel := getRoleLevel(req.Role)

		// 不能修改同级或上级成员的权限（除了自己）
		if targetLevel >= myLevel && user.ID != uint(userID) {
			c.JSON(http.StatusForbidden, response.Error(403, "权限不足：不能修改同级或上级成员的权限"))
			return
		}

		// 不能赋予比自己更高的权限
		if newLevel > myLevel {
			c.JSON(http.StatusForbidden, response.Error(403, "权限不足：不能赋予比自己更高的权限"))
			return
		}
	}

	// 规则：社长同时只能有一个，且只有学校管理员可以设置社长
	if req.Role == "leader" {
		if !authz.IsAdmin(user) {
			c.JSON(http.StatusForbidden, response.Error(403, "权限不足：只有学校管理员可以设置社长"))
			return
		}
		// 如果要设置为社长，先把该社团其他所有社长降级为member
		if err := store.DB().Model(&models.Membership{}).
			Where("club_id = ? AND role = ?", clubID, "leader").
			Update("role", "member").Error; err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(500, "更新失败"))
			return
		}
	}

	// 规则：如果目标用户当前是社长，只有学校管理员可以修改其权限（降级）
	if m.Role == "leader" && !authz.IsAdmin(user) {
		c.JSON(http.StatusForbidden, response.Error(403, "权限不足：只有学校管理员可以修改社长权限"))
		return
	}

	m.Role = req.Role
	if err := store.DB().Save(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "更新失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(m))
}
