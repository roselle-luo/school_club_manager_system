package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"web_server/db/models"
	"web_server/internal/store"
	"web_server/pkg/pagination"
	"web_server/pkg/response"

	"github.com/gin-gonic/gin"
)

// @Summary 公告列表（公开）
// @Tags 公共
// @Produce json
// @Param clubId query int false "社团ID"
// @Param keyword query string false "关键词"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Body
// @Router /public/announcements [get]
func ListPublicAnnouncements(c *gin.Context) {
	q := store.DB().Model(&models.Announcement{}).Where("scope = ?", "public")
	q = q.Preload("Club")
	if cid := c.Query("clubId"); cid != "" {
		if v, err := strconv.Atoi(cid); err == nil && v > 0 {
			q = q.Where("club_id = ?", v)
		}
	}
	if kw := strings.TrimSpace(c.Query("keyword")); kw != "" {
		like := "%%" + kw + "%%"
		q = q.Where("title LIKE ? OR content LIKE ?", like, like)
	}
	var list []models.Announcement
	pg := pagination.Get(c)
	info, err := pagination.Do(q.Order("id DESC"), pg, &list)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(map[string]any{"list": list, "pagination": info}))
}

// @Summary 活动列表（公开）
// @Tags 公共
// @Produce json
// @Param clubId query int false "社团ID"
// @Param keyword query string false "关键词"
// @Param start query string false "开始时间"
// @Param end query string false "结束时间"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Body
// @Router /public/activities [get]
func ListPublicActivities(c *gin.Context) {
	q := store.DB().Model(&models.Activity{}).Where("scope = ?", "public").Preload("Club")
	if cid := c.Query("clubId"); cid != "" {
		if v, err := strconv.Atoi(cid); err == nil && v > 0 {
			q = q.Where("club_id = ?", v)
		}
	}
	if kw := strings.TrimSpace(c.Query("keyword")); kw != "" {
		like := "%%" + kw + "%%"
		q = q.Where("subject LIKE ? OR place LIKE ?", like, like)
	}
	if start := strings.TrimSpace(c.Query("start")); start != "" {
		q = q.Where("(start_at IS NOT NULL AND start_at >= ?) OR (start_at IS NULL AND time >= ?)", start, start)
	}
	if end := strings.TrimSpace(c.Query("end")); end != "" {
		q = q.Where("(end_at IS NOT NULL AND end_at <= ?) OR (end_at IS NULL AND time <= ?)", end, end)
	}
	var list []models.Activity
	pg := pagination.Get(c)
	info, err := pagination.Do(q.Order("id DESC"), pg, &list)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(map[string]any{"list": list, "pagination": info}))
}

// @Summary 活动详情（公开）
// @Tags 公共
// @Produce json
// @Param activityId path int true "活动ID"
// @Success 200 {object} response.Body
// @Router /public/activities/{activityId} [get]
func GetPublicActivityDetail(c *gin.Context) {
	aidStr := c.Param("activityId")
	aid, err := strconv.Atoi(aidStr)
	if err != nil || aid <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	var a models.Activity
	if err := store.DB().Where("id = ?", aid).Preload("Club").First(&a).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "活动不存在"))
		return
	}
	c.JSON(http.StatusOK, response.Success(a))
}

// @Summary 社团详情（公开）
// @Tags 公共
// @Produce json
// @Param clubId path int true "社团ID"
// @Success 200 {object} response.Body
// @Router /public/clubs/{clubId} [get]
func GetClubDetailPublic(c *gin.Context) {
	clubIDStr := c.Param("clubId")
	clubID, err := strconv.Atoi(clubIDStr)
	if err != nil || clubID <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	var club models.Club
	if err := store.DB().Where("id = ?", clubID).Preload("Category").First(&club).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "社团不存在"))
		return
	}
	var leaders []models.Membership
	_ = store.DB().Where("club_id = ? AND role IN ?", club.ID, []string{"leader", "advisor"}).Preload("User").Find(&leaders)
	var acts []models.Activity
	_ = store.DB().Where("club_id = ? AND scope = ?", club.ID, "public").Order("id DESC").Limit(5).Find(&acts)
	// 统计成员数量（当前与历史）
	var currentCount int64
	_ = store.DB().Model(&models.Membership{}).
		Where("club_id = ? AND status = ?", club.ID, "approved").
		Count(&currentCount).Error
	var historyCount int64
	_ = store.DB().Model(&models.Membership{}).
		Where("club_id = ? AND status IN ?", club.ID, []string{"approved", "quit"}).
		Distinct("user_id").
		Count(&historyCount).Error
	res := map[string]any{
		"id":                   club.ID,
		"name":                 club.Name,
		"logo":                 club.Logo,
		"intro":                club.Intro,
		"category":             club.Category.Name,
		"leaders":              leaders,
		"activities":           acts,
		"member_count":         currentCount,
		"history_member_count": historyCount,
		"created_at":           club.CreatedAt,
	}
	c.JSON(http.StatusOK, response.Success(res))
}

// @Summary 社团类别列表（公开）
// @Tags 公共
// @Produce json
// @Success 200 {object} response.Body
// @Router /public/categories [get]
func ListCategories(c *gin.Context) {
	var list []models.ClubCategory
	q := store.DB().Model(&models.ClubCategory{}).Order("id DESC")
	pg := pagination.Get(c)
	info, err := pagination.Do(q, pg, &list)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(map[string]any{"list": list, "pagination": info}))
}
