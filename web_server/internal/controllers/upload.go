package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"web_server/config"
	"web_server/pkg/response"
)

// @Summary 上传图片
// @Tags 上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "图片文件"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /upload/image [post]
func UploadImage(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "缺少文件"))
		return
	}
	ext := strings.ToLower(filepath.Ext(f.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
	default:
		c.JSON(http.StatusBadRequest, response.Error(400, "不支持的图片类型"))
		return
	}

	cfg := config.Default()
	upRel := filepath.Join(cfg.Server.UploadDir)
	pubRel := filepath.Join(cfg.Server.PublicDir)
	_ = os.MkdirAll(filepath.Join(pubRel, upRel), 0755)
	name := time.Now().Format("20060102150405") + ext
	dst := filepath.Join(pubRel, upRel, name)
	if err := c.SaveUploadedFile(f, dst); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "保存失败"))
		return
	}
	url := "/static/" + filepath.ToSlash(filepath.Join(upRel, name))
	c.JSON(http.StatusOK, response.Success(map[string]string{"url": url}))
}
