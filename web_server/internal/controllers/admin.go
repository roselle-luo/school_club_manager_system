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

// @Summary 获取管理考勤列表（管理员/负责人）
// @Tags 管理员
// @Produce json
// @Param page query int false "页码"
// @Param size query int false "每页数量"
// @Param club_name query string false "社团名称"
// @Param user_name query string false "成员名字"
// @Param student_no query string false "学号"
// @Param date query string false "日期(YYYY-MM-DD)"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/attendance/list [get]
func ListManagedAttendance(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	offset := (page - 1) * size

	clubName := c.Query("club_name")
	clubIDStr := c.Query("club_id")
	userName := c.Query("user_name")
	studentNo := c.Query("student_no")
	dateStr := c.Query("date")

	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)

	db := store.DB().Model(&models.Attendance{})

	// 权限控制：如果不是管理员，只能查看自己负责的社团的考勤
	if !authz.IsAdmin(u) {
		var clubIDs []uint
		store.DB().Model(&models.Membership{}).
			Where("user_id = ? AND role IN ?", u.ID, []string{"leader", "advisor"}).
			Pluck("club_id", &clubIDs)

		if len(clubIDs) == 0 {
			// 如果没有管理的社团，直接返回空列表
			c.JSON(http.StatusOK, response.Success(map[string]any{
				"list":  []models.Attendance{},
				"total": 0,
			}))
			return
		}
		// 如果传入了club_id，必须检查该club_id是否在管理的社团列表中
		if clubIDStr != "" {
			cid, _ := strconv.Atoi(clubIDStr)
			allowed := false
			for _, id := range clubIDs {
				if id == uint(cid) {
					allowed = true
					break
				}
			}
			if !allowed {
				c.JSON(http.StatusForbidden, response.Error(403, "无权限查看该社团考勤"))
				return
			}
			db = db.Where("attendances.club_id = ?", cid)
		} else {
			db = db.Where("attendances.club_id IN ?", clubIDs)
		}
	} else {
		// 管理员：如果传入了club_id，则按club_id筛选
		if clubIDStr != "" {
			db = db.Where("attendances.club_id = ?", clubIDStr)
		}
	}

	if clubName != "" {
		db = db.Joins("JOIN clubs ON clubs.id = attendances.club_id").Where("clubs.name LIKE ?", "%"+clubName+"%")
	}
	if userName != "" || studentNo != "" {
		db = db.Joins("JOIN users ON users.id = attendances.user_id")
		if userName != "" {
			db = db.Where("users.name LIKE ?", "%"+userName+"%")
		}
		if studentNo != "" {
			db = db.Where("users.student_no LIKE ?", "%"+studentNo+"%")
		}
	}
	if dateStr != "" {
		// 筛选签到时间或签退时间匹配日期的记录
		// MySQL DATE() function
		db = db.Where("DATE(attendances.signin_at) = ? OR DATE(attendances.signout_at) = ?", dateStr, dateStr)
	}

	var total int64
	db.Count(&total)

	var list []models.Attendance
	if err := db.Preload("User").Preload("Club").Preload("Activity").
		Order("attendances.id desc").Offset(offset).Limit(size).Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(map[string]any{
		"list":  list,
		"total": total,
	}))
}
