package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"web_server/db/models"
	"web_server/internal/authz"
	"web_server/internal/store"
	"web_server/pkg/pagination"
	"web_server/pkg/response"

	"github.com/gin-gonic/gin"
)

// @Summary 获取当前用户管理的社团
// @Tags 负责人
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/clubs [get]
func GetUserLeaderClubs(c *gin.Context) {
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)

	type ClubWithRole struct {
		models.Club
		Role string `json:"role"`
	}
	var result []ClubWithRole

	if authz.IsAdmin(u) {
		var clubs []models.Club
		if err := store.DB().Preload("Category").Find(&clubs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
			return
		}
		for _, club := range clubs {
			result = append(result, ClubWithRole{Club: club, Role: "admin"})
		}
	} else {
		// 查找担任leader或advisor的社团
		var memberships []models.Membership
		store.DB().Where("user_id = ? AND role IN ?", u.ID, []string{"leader", "advisor"}).Find(&memberships)

		clubIDs := make([]uint, 0)
		roleMap := make(map[uint]string)
		for _, m := range memberships {
			clubIDs = append(clubIDs, m.ClubID)
			roleMap[m.ClubID] = m.Role
		}

		if len(clubIDs) > 0 {
			var clubs []models.Club
			store.DB().Preload("Category").Where("id IN ?", clubIDs).Find(&clubs)
			for _, club := range clubs {
				result = append(result, ClubWithRole{Club: club, Role: roleMap[club.ID]})
			}
		}
	}
	c.JSON(http.StatusOK, response.Success(result))
}

// @Summary 获取社团负责人列表
// @Tags 负责人
// @Produce json
// @Param clubId path int true "社团ID"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/clubs/{clubId}/leaders [get]
func GetClubLeaders(c *gin.Context) {
	clubIDStr := c.Param("clubId")
	clubID, _ := strconv.Atoi(clubIDStr)

	var members []models.Membership
	if err := store.DB().Preload("User").Where("club_id = ? AND role IN ?", clubID, []string{"leader", "advisor"}).Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(members))
}

type SetRoleReq struct {
	Role string `json:"role" binding:"required"` // member, leader, advisor
}

// @Summary 负责人设定成员角色
// @Tags 负责人
// @Accept json
// @Produce json
// @Param clubId path int true "社团ID"
// @Param userId path int true "成员用户ID"
// @Param payload body SetRoleReq true "角色信息"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/clubs/{clubId}/members/{userId}/role [put]
func SetMemberRoleByLeader(c *gin.Context) {
	clubIDStr := c.Param("clubId")
	userIDStr := c.Param("userId")
	clubID, _ := strconv.Atoi(clubIDStr)
	userID, _ := strconv.Atoi(userIDStr)

	cu, _ := c.Get("currentUser")
	user := cu.(*models.User)

	var req SetRoleReq
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

	m.Role = req.Role
	if err := store.DB().Save(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "更新失败"))
		return
	}
	RecordLog(user.ID, user.Name, "修改权限", fmt.Sprintf("修改用户 %d 角色为 %s", userID, req.Role), uint(clubID))
	c.JSON(http.StatusOK, response.Success(m))
}

// @Summary 待审批入会列表
// @Tags 成员
// @Produce json
// @Param clubId path int true "社团ID"
// @Security Bearer
// @Success 200 {object} response.Body
// ListPendingMembershipsOld (Renamed due to duplication)
// @Router /leader/clubs/{clubId}/memberships/pending/old [get]
func ListPendingMembershipsOld(c *gin.Context) {
	clubID, _ := strconv.Atoi(c.Param("clubId"))
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	if !(authz.IsAdmin(u) || authz.IsClubLeader(u.ID, uint(clubID))) {
		c.JSON(http.StatusForbidden, response.Error(403, "无权限"))
		return
	}
	var list []models.Membership
	if err := store.DB().Preload("User").Where("club_id = ? AND status = ?", clubID, "pending").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(list))
}

// @Summary 审批通过成员（负责人）
// @Tags 成员
// @Accept json
// @Produce json
// @Param clubId path int true "社团ID"
// @Param id path int true "成员关系ID"
// @Security Bearer
// @Success 200 {object} response.Body
// ApproveMembershipOld (Renamed due to duplication)
// @Router /leader/clubs/{clubId}/memberships/{id}/approve/old [post]
func ApproveMembershipOld(c *gin.Context) {
	clubID, _ := strconv.Atoi(c.Param("clubId"))
	id, _ := strconv.Atoi(c.Param("id"))
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	if !(authz.IsAdmin(u) || authz.IsClubLeader(u.ID, uint(clubID))) {
		c.JSON(http.StatusForbidden, response.Error(403, "无权限"))
		return
	}
	var m models.Membership
	if err := store.DB().Where("id = ? AND club_id = ?", id, clubID).First(&m).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "成员不存在"))
		return
	}
	m.Status = "approved"
	if err := store.DB().Save(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "更新失败"))
		return
	}
	RecordLog(u.ID, u.Name, "审批申请", fmt.Sprintf("批准成员 %d 加入社团", m.UserID), uint(clubID))
	c.JSON(http.StatusOK, response.Success(m))
}

// @Summary 拒绝成员加入（负责人）
// @Tags 成员
// @Accept json
// @Produce json
// @Param clubId path int true "社团ID"
// @Param id path int true "成员关系ID"
// @Security Bearer
// @Success 200 {object} response.Body
// RejectMembershipOld (Renamed due to duplication)
// @Router /leader/clubs/{clubId}/memberships/{id}/reject/old [post]
func RejectMembershipOld(c *gin.Context) {
	clubID, _ := strconv.Atoi(c.Param("clubId"))
	id, _ := strconv.Atoi(c.Param("id"))
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	if !(authz.IsAdmin(u) || authz.IsClubLeader(u.ID, uint(clubID))) {
		c.JSON(http.StatusForbidden, response.Error(403, "无权限"))
		return
	}
	var m models.Membership
	if err := store.DB().Where("id = ? AND club_id = ?", id, clubID).First(&m).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "成员不存在"))
		return
	}
	m.Status = "rejected"
	if err := store.DB().Save(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "更新失败"))
		return
	}
	RecordLog(u.ID, u.Name, "审批申请", fmt.Sprintf("拒绝成员 %d 加入社团", m.UserID), uint(clubID))
	c.JSON(http.StatusOK, response.Success(m))
}

// @Summary 社团成员列表（负责人）
// @Tags 成员
// @Produce json
// @Param clubId path int true "社团ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/clubs/{clubId}/members [get]
func ListClubMembersLegacy(c *gin.Context) {
	clubID, _ := strconv.Atoi(c.Param("clubId"))
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	if !(authz.IsAdmin(u) || authz.IsClubLeader(u.ID, uint(clubID))) {
		c.JSON(http.StatusForbidden, response.Error(403, "无权限"))
		return
	}
	var list []models.Membership
	q := store.DB().Preload("User").Where("club_id = ? AND status = ?", clubID, "approved")
	pg := pagination.Get(c)
	info, err := pagination.Do(q, pg, &list)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(map[string]any{"list": list, "pagination": info}))
}

// @Summary 解散社团
// @Tags 负责人
// @Produce json
// @Param clubId path int true "社团ID"
// @Security Bearer
// @Success 200 {object} response.Body
// DissolveClubOld (Renamed due to duplication)
// @Router /leader/clubs/{clubId}/old [delete]
func DissolveClubOld(c *gin.Context) {
	clubID, _ := strconv.Atoi(c.Param("clubId"))
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	if !(authz.IsAdmin(u) || authz.IsClubLeader(u.ID, uint(clubID))) {
		c.JSON(http.StatusForbidden, response.Error(403, "无权限"))
		return
	}
	if err := store.DB().Delete(&models.Club{}, clubID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "删除失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}

// @Summary 踢出成员
// @Tags 负责人
// @Produce json
// @Param clubId path int true "社团ID"
// @Param userId path int true "用户ID"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/clubs/{clubId}/members/{userId} [delete]
func KickMember(c *gin.Context) {
	clubID, _ := strconv.Atoi(c.Param("clubId"))
	userID, _ := strconv.Atoi(c.Param("userId"))

	cu, _ := c.Get("currentUser")
	caller := cu.(*models.User)

	// Get target membership
	var targetM models.Membership
	if err := store.DB().Where("user_id = ? AND club_id = ?", userID, clubID).First(&targetM).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "成员不存在该社团"))
		return
	}

	// Permission Logic
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

	targetLevel := getRoleLevel(targetM.Role)

	// Special rule: Leader cannot be kicked by anyone (even Admin) via this interface
	if targetM.Role == "leader" {
		c.JSON(http.StatusForbidden, response.Error(403, "无法踢出社长"))
		return
	}

	if authz.IsAdmin(caller) {
		// Admin can kick anyone except leader (handled above)
	} else {
		// Club Leader/Advisor checks
		var callerM models.Membership
		if err := store.DB().Where("user_id = ? AND club_id = ?", caller.ID, clubID).First(&callerM).Error; err != nil {
			c.JSON(http.StatusForbidden, response.Error(403, "无权限"))
			return
		}

		myLevel := getRoleLevel(callerM.Role)

		// Must have higher level than target
		if myLevel <= targetLevel {
			c.JSON(http.StatusForbidden, response.Error(403, "权限不足：只能踢出权限低于自己的成员"))
			return
		}
	}

	// Delete membership
	if err := store.DB().Delete(&targetM).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "操作失败"))
		return
	}

	RecordLog(caller.ID, caller.Name, "修改权限", fmt.Sprintf("将成员 %d 踢出社团", userID), uint(clubID))
	c.JSON(http.StatusOK, response.Success(nil))
}
