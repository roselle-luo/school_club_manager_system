package controllers

import (
	"math"
	"net/http"
	"strconv"
	"time"
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
	// 必须先报名
	var reg models.ActivityParticipant
	if err := store.DB().Where("user_id = ? AND activity_id = ? AND club_id = ? AND status = ?", u.ID, activityID, act.ClubID, "confirmed").First(&reg).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "需先报名该活动"))
		return
	}
	// 检查是否已有未签退的签到记录
	var latest models.Attendance
	if err := store.DB().Where("user_id = ? AND club_id = ? AND activity_id = ? AND signout_at IS NULL", u.ID, act.ClubID, activityID).Order("id DESC").First(&latest).Error; err == nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "已签到，未签退"))
		return
	}
	now := time.Now()
	aid := uint(activityID)
	att := models.Attendance{UserID: u.ID, ActivityID: &aid, ClubID: act.ClubID, SigninAt: &now}
	if err := store.DB().Create(&att).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "签到失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(att))
}

// @Summary 社团签到（与活动无关）
// @Tags 考勤
// @Produce json
// @Param clubId path int true "社团ID"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /member/clubs/{clubId}/signin [post]
func ClubSignIn(c *gin.Context) {
	clubIDStr := c.Param("clubId")
	clubID, err := strconv.Atoi(clubIDStr)
	if err != nil || clubID <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	var cl models.Club
	if err := store.DB().Where("id = ?", clubID).First(&cl).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "社团不存在"))
		return
	}
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	if !authz.IsClubMember(u.ID, cl.ID) {
		c.JSON(http.StatusForbidden, response.Error(403, "非社团成员"))
		return
	}
	var latest models.Attendance
	if err := store.DB().Where("user_id = ? AND club_id = ? AND signout_at IS NULL", u.ID, clubID).Order("id DESC").First(&latest).Error; err == nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "已签到，未签退"))
		return
	}
	now := time.Now()
	att := models.Attendance{UserID: u.ID, ActivityID: nil, ClubID: uint(clubID), SigninAt: &now}
	if err := store.DB().Create(&att).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "签到失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(att))
}

// @Summary 社团签退（与活动无关）
// @Tags 考勤
// @Produce json
// @Param clubId path int true "社团ID"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /member/clubs/{clubId}/signout [post]
func ClubSignOut(c *gin.Context) {
	clubIDStr := c.Param("clubId")
	clubID, err := strconv.Atoi(clubIDStr)
	if err != nil || clubID <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	var cl models.Club
	if err := store.DB().Where("id = ?", clubID).First(&cl).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "社团不存在"))
		return
	}
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	if !authz.IsClubMember(u.ID, cl.ID) {
		c.JSON(http.StatusForbidden, response.Error(403, "非社团成员"))
		return
	}
	var latest models.Attendance
	if err := store.DB().Where("user_id = ? AND club_id = ? AND signout_at IS NULL", u.ID, clubID).Order("id DESC").First(&latest).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "未找到签到记录"))
		return
	}
	now := time.Now()
	if latest.SigninAt != nil {
		d := now.Sub(*latest.SigninAt)
		if d < time.Minute {
			if err := store.DB().Delete(&latest).Error; err != nil {
				c.JSON(http.StatusInternalServerError, response.Error(500, "签退失败"))
				return
			}
			c.JSON(http.StatusOK, response.Success(nil))
			return
		}
		latest.DurationMinutes = int(d.Minutes())
		hours := d.Hours()
		latest.DurationHours = math.Round(hours*100) / 100
	}
	latest.SignoutAt = &now
	if err := store.DB().Save(&latest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "签退失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(latest))
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
	// 查找最近一次未签退的签到记录并更新为签退
	var latest models.Attendance
	if err := store.DB().Where("user_id = ? AND club_id = ? AND activity_id = ? AND signout_at IS NULL", u.ID, act.ClubID, activityID).Order("id DESC").First(&latest).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "未找到签到记录"))
		return
	}
	now := time.Now()
	if latest.SigninAt != nil {
		d := now.Sub(*latest.SigninAt)
		if d < time.Minute {
			if err := store.DB().Delete(&latest).Error; err != nil {
				c.JSON(http.StatusInternalServerError, response.Error(500, "签退失败"))
				return
			}
			c.JSON(http.StatusOK, response.Success(nil))
			return
		}
		latest.DurationMinutes = int(d.Minutes())
		hours := d.Hours()
		latest.DurationHours = math.Round(hours*100) / 100
	}
	latest.SignoutAt = &now
	if err := store.DB().Save(&latest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "签退失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(latest))
}

// @Summary 报名参加活动
// @Tags 活动
// @Produce json
// @Param activityId path int true "活动ID"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /member/activities/{activityId}/register [post]
func RegisterActivity(c *gin.Context) {
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
	// 已报名则返回成功
	var exist models.ActivityParticipant
	if err := store.DB().Where("user_id = ? AND activity_id = ? AND club_id = ?", u.ID, activityID, act.ClubID).First(&exist).Error; err == nil {
		if exist.Status != "confirmed" {
			exist.Status = "confirmed"
			_ = store.DB().Save(&exist).Error
		}
		c.JSON(http.StatusOK, response.Success(exist))
		return
	}
	// 人数限制
	if act.MaxParticipants > 0 {
		var cnt int64
		_ = store.DB().Model(&models.ActivityParticipant{}).Where("activity_id = ? AND status = ?", activityID, "confirmed").Count(&cnt)
		if int(cnt) >= act.MaxParticipants {
			c.JSON(http.StatusBadRequest, response.Error(400, "报名人数已满"))
			return
		}
	}
	reg := models.ActivityParticipant{UserID: u.ID, ActivityID: uint(activityID), ClubID: act.ClubID, Status: "confirmed"}
	if err := store.DB().Create(&reg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "报名失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(reg))
}

// @Summary 报名状态查询
// @Tags 活动
// @Produce json
// @Param activityId path int true "活动ID"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /member/activities/{activityId}/register [get]
func GetRegisterStatus(c *gin.Context) {
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
	var exist models.ActivityParticipant
	err = store.DB().Where("user_id = ? AND activity_id = ? AND club_id = ? AND status = ?", u.ID, activityID, act.ClubID, "confirmed").First(&exist).Error
	c.JSON(http.StatusOK, response.Success(map[string]any{
		"registered": err == nil,
	}))
}

// @Summary 取消报名
// @Tags 活动
// @Produce json
// @Param activityId path int true "活动ID"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /member/activities/{activityId}/register [delete]
func CancelRegisterActivity(c *gin.Context) {
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
	var exist models.ActivityParticipant
	if err := store.DB().Where("user_id = ? AND activity_id = ? AND club_id = ? AND status = ?", u.ID, activityID, act.ClubID, "confirmed").First(&exist).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "未报名"))
		return
	}
	exist.Status = "cancelled"
	if err := store.DB().Save(&exist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "取消失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
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
	q := store.DB().Model(&models.Attendance{}).Where("user_id = ? AND club_id = ?", u.ID, clubID).Preload("Activity").Order("id DESC")
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

// @Summary 强制签退（负责人）
// @Tags 考勤
// @Produce json
// @Param id path int true "考勤记录ID"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/attendance/{id}/signout [post]
func ForceSignOut(c *gin.Context) {
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
	if att.SignoutAt != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "该记录已签退"))
		return
	}
	now := time.Now()
	if att.SigninAt != nil {
		d := now.Sub(*att.SigninAt)
		att.DurationMinutes = int(d.Minutes())
		hours := d.Hours()
		att.DurationHours = math.Round(hours*100) / 100
	}
	att.SignoutAt = &now
	if err := store.DB().Save(&att).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "强制签退失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(att))
}
