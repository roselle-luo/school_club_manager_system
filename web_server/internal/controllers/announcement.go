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

type AnnouncementReq struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Scope   string `json:"scope"`
}

// @Summary 获取社团公告列表
// @Tags 负责人
// @Produce json
// @Param clubId path int true "社团ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/clubs/{clubId}/announcements [get]
func ListClubAnnouncements(c *gin.Context) {
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

	var list []models.Announcement
	q := store.DB().Model(&models.Announcement{}).Where("club_id = ?", clubID).Order("id DESC")
	pg := pagination.Get(c)
	info, err := pagination.Do(q, pg, &list)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(map[string]any{"list": list, "pagination": info}))
}

// @Summary 创建社团公告
// @Tags 负责人
// @Accept json
// @Produce json
// @Param clubId path int true "社团ID"
// @Param payload body AnnouncementReq true "公告内容"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/clubs/{clubId}/announcements [post]
func CreateAnnouncement(c *gin.Context) {
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

	var req AnnouncementReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}

	ann := models.Announcement{
		Title:   req.Title,
		Content: req.Content,
		Scope:   req.Scope,
		ClubID:  uint(clubID),
	}
	if ann.Scope == "" {
		ann.Scope = "public"
	}

	if err := store.DB().Create(&ann).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "创建失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(ann))
}

// @Summary 更新社团公告
// @Tags 负责人
// @Accept json
// @Produce json
// @Param clubId path int true "社团ID"
// @Param id path int true "公告ID"
// @Param payload body AnnouncementReq true "公告内容"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/clubs/{clubId}/announcements/{id} [put]
func UpdateAnnouncement(c *gin.Context) {
	clubIDStr := c.Param("clubId")
	idStr := c.Param("id")
	clubID, err1 := strconv.Atoi(clubIDStr)
	id, err2 := strconv.Atoi(idStr)
	if err1 != nil || err2 != nil || clubID <= 0 || id <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	if !(authz.IsAdmin(u) || authz.IsClubLeader(u.ID, uint(clubID))) {
		c.JSON(http.StatusForbidden, response.Error(403, "无权限"))
		return
	}

	var req AnnouncementReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}

	var ann models.Announcement
	if err := store.DB().Where("id = ? AND club_id = ?", id, clubID).First(&ann).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "公告不存在"))
		return
	}

	updates := map[string]any{
		"title":   req.Title,
		"content": req.Content,
	}
	if req.Scope != "" {
		updates["scope"] = req.Scope
	}

	if err := store.DB().Model(&ann).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "更新失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(ann))
}

// @Summary 删除社团公告
// @Tags 负责人
// @Produce json
// @Param clubId path int true "社团ID"
// @Param id path int true "公告ID"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/clubs/{clubId}/announcements/{id} [delete]
func DeleteAnnouncement(c *gin.Context) {
	clubIDStr := c.Param("clubId")
	idStr := c.Param("id")
	clubID, err1 := strconv.Atoi(clubIDStr)
	id, err2 := strconv.Atoi(idStr)
	if err1 != nil || err2 != nil || clubID <= 0 || id <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	if !(authz.IsAdmin(u) || authz.IsClubLeader(u.ID, uint(clubID))) {
		c.JSON(http.StatusForbidden, response.Error(403, "无权限"))
		return
	}

	if err := store.DB().Where("id = ? AND club_id = ?", id, clubID).Delete(&models.Announcement{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "删除失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}
