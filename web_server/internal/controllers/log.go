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

// @Summary 获取社团操作日志
// @Tags 负责人
// @Produce json
// @Param clubId path int true "社团ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param actionType query string false "操作类型"
// @Param operatorName query string false "操作人姓名"
// @Security Bearer
// @Success 200 {object} response.Body
// ListOperationLogsOld (Renamed due to duplication)
// @Router /leader/clubs/{clubId}/logs/old [get]
func ListOperationLogsOld(c *gin.Context) {
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

	q := store.DB().Model(&models.OperationLog{}).Where("club_id = ?", clubID).Order("id DESC")

	if actionType := c.Query("actionType"); actionType != "" {
		q = q.Where("action_type = ?", actionType)
	}
	if operatorName := c.Query("operatorName"); operatorName != "" {
		q = q.Where("operator_name LIKE ?", "%"+operatorName+"%")
	}

	var list []models.OperationLog
	pg := pagination.Get(c)
	info, err := pagination.Do(q, pg, &list)
	if err != nil {
		fmt.Printf("ListOperationLogs error: %v\n", err)
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(map[string]any{"list": list, "pagination": info}))
}
