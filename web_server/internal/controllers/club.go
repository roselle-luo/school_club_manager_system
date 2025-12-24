package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"web_server/db/models"
	"web_server/internal/authz"
	"web_server/internal/store"
	"web_server/pkg/pagination"
	"web_server/pkg/response"

	"github.com/gin-gonic/gin"
)

type ClubItem struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Logo        string         `json:"logo"`
	Intro       string         `json:"intro"`
	Category    string         `json:"category"`
	LeaderName  string         `json:"leader_name"`
	LeaderPhone string         `json:"leader_phone"`
	Activities  []ActivityItem `json:"activities"`
}

type ActivityItem struct {
	ID      uint   `json:"id"`
	Subject string `json:"subject"`
	Time    string `json:"time"`
	Place   string `json:"place"`
}

// @Summary 社团列表（含类别、代表活动、负责人联系方式）
// @Tags 公共
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Body
// @Router /public/clubs [get]
func ListClubs(c *gin.Context) {
	var clubs []models.Club
	q := store.DB().Model(&models.Club{}).Preload("Category")
	if cid := c.Query("categoryId"); cid != "" {
		if v, err := strconv.Atoi(cid); err == nil && v > 0 {
			q = q.Where("category_id = ?", v)
		}
	}
	if kw := c.Query("keyword"); kw != "" {
		like := "%%" + kw + "%%"
		q = q.Where("name LIKE ? OR intro LIKE ?", like, like)
	}
	pg := pagination.Get(c)
	var totalInfo pagination.Info
	var err error
	totalInfo, err = pagination.Do(q.Order("id DESC"), pg, &clubs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
		return
	}
	items := make([]ClubItem, 0, len(clubs))
	for _, cl := range clubs {
		var leader models.Membership
		_ = store.DB().Where("club_id = ? AND role IN ?", cl.ID, []string{"leader", "advisor"}).Preload("User").Order("id ASC").First(&leader)
		var acts []models.Activity
		_ = store.DB().Where("club_id = ? AND scope = ?", cl.ID, "public").Order("id DESC").Limit(3).Find(&acts)
		ai := make([]ActivityItem, 0, len(acts))
		for _, a := range acts {
			ai = append(ai, ActivityItem{ID: a.ID, Subject: a.Subject, Time: a.Time, Place: a.Place})
		}
		items = append(items, ClubItem{
			ID:          cl.ID,
			Name:        cl.Name,
			Logo:        cl.Logo,
			Intro:       cl.Intro,
			Category:    cl.Category.Name,
			LeaderName:  leader.User.Name,
			LeaderPhone: leader.User.Phone,
			Activities:  ai,
		})
	}
	c.JSON(http.StatusOK, response.Success(map[string]any{"list": items, "pagination": totalInfo}))
}

// @Summary 申请加入社团
// @Tags 学生
// @Produce json
// @Param clubId path int true "社团ID"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /student/clubs/{clubId}/apply [post]
func ApplyJoinClub(c *gin.Context) {
	clubID := c.Param("clubId")
	var cl models.Club
	if err := store.DB().Where("id = ?", clubID).First(&cl).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "社团不存在"))
		return
	}
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	var m models.Membership
	if err := store.DB().Where("user_id = ? AND club_id = ?", u.ID, cl.ID).First(&m).Error; err == nil {
		m.Status = "pending"
		m.Role = "member"
		if err := store.DB().Save(&m).Error; err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(500, "申请失败"))
			return
		}
		c.JSON(http.StatusOK, response.Success(m))
		return
	}
	m = models.Membership{UserID: u.ID, ClubID: cl.ID, Status: "pending", Role: "member"}
	if err := store.DB().Create(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "申请失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(m))
}

// @Summary 我的社团关系
// @Tags 学生
// @Produce json
// @Param status query string false "状态: pending/approved/rejected/quit"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /student/memberships/my [get]
func MyMemberships(c *gin.Context) {
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	q := store.DB().Model(&models.Membership{}).Where("user_id = ?", u.ID).Preload("Club").Order("id DESC")
	if st := c.Query("status"); st != "" {
		q = q.Where("status = ?", st)
	}
	var list []models.Membership
	pg := pagination.Get(c)
	info, err := pagination.Do(q, pg, &list)
	if err != nil {
		log.Printf("MyMemberships query error: user_id=%d status=%s err=%v", u.ID, c.Query("status"), err)
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(map[string]any{"list": list, "pagination": info}))
}

// @Summary 退出社团
// @Tags 学生
// @Produce json
// @Param clubId path int true "社团ID"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /student/clubs/{clubId}/exit [post]
func ExitClub(c *gin.Context) {
	clubIDStr := c.Param("clubId")
	clubID, err := strconv.Atoi(clubIDStr)
	if err != nil || clubID <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	var m models.Membership
	if err := store.DB().Where("user_id = ? AND club_id = ?", u.ID, clubID).First(&m).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "未加入该社团"))
		return
	}
	m.Status = "quit"
	if err := store.DB().Save(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "退出失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(m))
}

// @Summary 待审批入会列表（负责人）
// @Tags 成员
// @Produce json
// @Param clubId path int true "社团ID"
// @Param keyword query string false "关键词：姓名/学号"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/clubs/{clubId}/memberships [get]
func ListPendingMemberships(c *gin.Context) {
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
	var list []models.Membership
	q := store.DB().Model(&models.Membership{}).Where("club_id = ? AND status = ?", clubID, "pending").Preload("User").Order("id DESC")
	if kw := c.Query("keyword"); kw != "" {
		like := "%%" + kw + "%%"
		q = q.Where("user_id IN (SELECT id FROM users WHERE name LIKE ? OR student_no LIKE ?)", like, like)
	}
	pg := pagination.Get(c)
	info, err := pagination.Do(q, pg, &list)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(map[string]any{"list": list, "pagination": info}))
}

// @Summary 审批通过成员（负责人）
// @Tags 成员
// @Accept json
// @Produce json
// @Param clubId path int true "社团ID"
// @Param id path int true "成员关系ID"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/clubs/{clubId}/memberships/{id}/approve [post]
func ApproveMembership(c *gin.Context) {
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

// @Summary 审批驳回成员（负责人）
// @Tags 成员
// @Accept json
// @Produce json
// @Param clubId path int true "社团ID"
// @Param id path int true "成员关系ID"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/clubs/{clubId}/memberships/{id}/reject [post]
func RejectMembership(c *gin.Context) {
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
	RecordLog(u.ID, u.Name, "审批申请", fmt.Sprintf("驳回成员 %d 加入社团", m.UserID), uint(clubID))
	c.JSON(http.StatusOK, response.Success(m))
}

// @Summary 社团成员列表（负责人）
// @Tags 成员
// @Produce json
// @Param clubId path int true "社团ID"
// @Param role query string false "角色: member/leader/advisor"
// @Param status query string false "状态: pending/approved/rejected/quit"
// @Param keyword query string false "关键词：姓名/学号/手机号"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/clubs/{clubId}/members/users [get]
func ListClubMembers(c *gin.Context) {
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
	q := store.DB().Model(&models.Membership{}).Where("club_id = ?", clubID).Preload("User").Order("id DESC")
	if role := c.Query("role"); role != "" {
		q = q.Where("role = ?", role)
	}
	if status := c.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	if kw := c.Query("keyword"); kw != "" {
		like := "%%" + kw + "%%"
		q = q.Where("user_id IN (SELECT id FROM users WHERE name LIKE ? OR student_no LIKE ? OR phone LIKE ?)", like, like, like)
	}
	var list []models.Membership
	pg := pagination.Get(c)
	info, err := pagination.Do(q, pg, &list)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(map[string]any{"list": list, "pagination": info}))
}

// @Summary 解散社团
// @Tags 管理员
// @Produce json
// @Param clubId path int true "社团ID"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /admin/clubs/{clubId} [delete]
func DissolveClub(c *gin.Context) {
	clubIDStr := c.Param("clubId")
	clubID, err := strconv.Atoi(clubIDStr)
	if err != nil || clubID <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	if !authz.IsAdmin(u) {
		c.JSON(http.StatusForbidden, response.Error(403, "权限不足：只有学校管理员可以解散社团"))
		return
	}

	tx := store.DB().Begin()
	// 删除社团成员
	if err := tx.Where("club_id = ?", clubID).Delete(&models.Membership{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, response.Error(500, "解散失败：删除成员失败"))
		return
	}
	// 删除社团活动
	if err := tx.Where("club_id = ?", clubID).Delete(&models.Activity{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, response.Error(500, "解散失败：删除活动失败"))
		return
	}
	// 删除社团
	if err := tx.Delete(&models.Club{}, clubID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, response.Error(500, "解散失败"))
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, response.Success(nil))
}
