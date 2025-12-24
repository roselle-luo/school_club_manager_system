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

// RecordLog 记录操作日志
// operatorID: 操作者ID
// operatorName: 操作者姓名
// actionType: 操作类型 (e.g., "审批申请", "修改打卡", "修改权限", "发布公告")
// content: 操作详情
// clubID: 相关社团ID (0表示系统级或无关社团)
func RecordLog(operatorID uint, operatorName string, actionType string, content string, clubID uint) {
	log := models.OperationLog{
		OperatorID:   operatorID,
		OperatorName: operatorName,
		ActionType:   actionType,
		Content:      content,
		ClubID:       clubID,
	}
	// 异步记录，避免阻塞主流程，但要注意错误处理（这里简单打印）
	go func() {
		if err := store.DB().Create(&log).Error; err != nil {
			fmt.Printf("Failed to record log: %v\n", err)
		}
	}()
}

// @Summary 获取操作日志列表
// @Tags 日志
// @Produce json
// @Param clubId path int true "社团ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param actionType query string false "操作类型"
// @Param operatorName query string false "操作者姓名"
// @Param startDate query string false "开始日期(YYYY-MM-DD)"
// @Param endDate query string false "结束日期(YYYY-MM-DD)"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /leader/clubs/{clubId}/logs [get]
func ListOperationLogs(c *gin.Context) {
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

	var list []models.OperationLog
	q := store.DB().Model(&models.OperationLog{}).Where("club_id = ?", clubID)

	if actionType := c.Query("actionType"); actionType != "" {
		q = q.Where("action_type = ?", actionType)
	}
	if operatorName := c.Query("operatorName"); operatorName != "" {
		q = q.Where("operator_name LIKE ?", "%"+operatorName+"%")
	}
	if startDate := c.Query("startDate"); startDate != "" {
		q = q.Where("created_at >= ?", startDate+" 00:00:00")
	}
	if endDate := c.Query("endDate"); endDate != "" {
		q = q.Where("created_at <= ?", endDate+" 23:59:59")
	}

	pg := pagination.Get(c)
	info, err := pagination.Do(q.Order("id DESC"), pg, &list)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(map[string]any{"list": list, "pagination": info}))
}
