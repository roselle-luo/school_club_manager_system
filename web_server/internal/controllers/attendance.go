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

// @Summary 会员签到
// @Tags 考勤
// @Produce json
// @Param activityId path int true "活动ID"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /member/activities/{activityId}/signin [post]
func SignIn(c *gin.Context) {
	activityIDStr := c.Param("activityId")
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil || activityID <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	var act models.Activity
	if err := store.DB().Where("id = ?", activityID).First(&act).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "活动不存在"))
		return
	}
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	if !authz.IsClubMember(u.ID, act.ClubID) {
		c.JSON(http.StatusForbidden, response.Error(403, "非社团成员"))
		return
	}
	att := models.Attendance{UserID: u.ID, ActivityID: uint(activityID), ClubID: act.ClubID, Type: "signin"}
	if err := store.DB().Create(&att).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "签到失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(att))
}

// @Summary 会员签退
// @Tags 考勤
// @Produce json
// @Param activityId path int true "活动ID"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /member/activities/{activityId}/signout [post]
func SignOut(c *gin.Context) {
	activityIDStr := c.Param("activityId")
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil || activityID <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	var act models.Activity
	if err := store.DB().Where("id = ?", activityID).First(&act).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "活动不存在"))
		return
	}
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	if !authz.IsClubMember(u.ID, act.ClubID) {
		c.JSON(http.StatusForbidden, response.Error(403, "非社团成员"))
		return
	}
	att := models.Attendance{UserID: u.ID, ActivityID: uint(activityID), ClubID: act.ClubID, Type: "signout"}
	if err := store.DB().Create(&att).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "签退失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(att))
}

// @Summary 我的考勤记录（按社团）
// @Tags 考勤
// @Produce json
// @Param clubId query int true "社团ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /member/attendance/my [get]
func MyAttendance(c *gin.Context) {
	clubIDStr := c.Query("clubId")
	clubID, err := strconv.Atoi(clubIDStr)
	if err != nil || clubID <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	if !authz.IsClubMember(u.ID, uint(clubID)) {
		c.JSON(http.StatusForbidden, response.Error(403, "非社团成员"))
		return
	}
	var list []models.Attendance
	q := store.DB().Model(&models.Attendance{}).Where("user_id = ? AND club_id = ?", u.ID, clubID).Order("id DESC")
	pg := pagination.Get(c)
	info, err := pagination.Do(q, pg, &list)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(map[string]any{"list": list, "pagination": info}))
}

// @Summary 社团考勤检索（负责人）
// @Tags 考勤
// @Produce json
// @Param clubId path int true "社团ID"
// @Param userId query int false "用户ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/clubs/{clubId}/attendance [get]
func ClubAttendance(c *gin.Context) {
	clubIDStr := c.Param("clubId")
	clubID, err := strconv.Atoi(clubIDStr)
	if err != nil || clubID <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	if !(authz.IsAdmin(u) || authz.IsClubLeader(u.ID, uint(clubID))) {
		c.JSON(http.StatusForbidden, response.Error(403, "无权限"))
		return
	}
	q := store.DB().Where("club_id = ?", clubID)
	if uidStr := c.Query("userId"); uidStr != "" {
		if uid, e := strconv.Atoi(uidStr); e == nil && uid > 0 {
			q = q.Where("user_id = ?", uid)
		}
	}
	var list []models.Attendance
	pg := pagination.Get(c)
	info, err := pagination.Do(q.Order("id DESC"), pg, &list)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(map[string]any{"list": list, "pagination": info}))
}

// @Summary 删除异常考勤记录（负责人）
// @Tags 考勤
// @Produce json
// @Param id path int true "考勤记录ID"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/attendance/{id} [delete]
func DeleteAttendance(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	var att models.Attendance
	if err := store.DB().Where("id = ?", id).First(&att).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "不存在"))
		return
	}
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	if !(authz.IsAdmin(u) || authz.IsClubLeader(u.ID, att.ClubID)) {
		c.JSON(http.StatusForbidden, response.Error(403, "无权限"))
		return
	}
	if err := store.DB().Delete(&att).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "删除失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}
