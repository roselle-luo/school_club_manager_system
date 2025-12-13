package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web_server/config"
	"web_server/db/models"
	"web_server/internal/store"
	"web_server/pkg/jwt"
	"web_server/pkg/password"
	"web_server/pkg/response"
)

type LoginReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary 登录
// @Tags 公共
// @Accept json
// @Produce json
// @Param payload body LoginReq true "登录参数"
// @Success 200 {object} response.Body
// @Router /public/login [post]
func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	var u models.User
	if err := store.DB().Where("account = ?", req.Account).Preload("Role").First(&u).Error; err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(401, "账号或密码错误"))
		return
	}
	if !password.Compare(u.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, response.Error(401, "账号或密码错误"))
		return
	}
	token, err := jwt.Sign(config.Default().JWT.Secret, u.ID, u.Role.Code, config.Default().JWT.Expires)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "登录失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(map[string]any{"token": token}))
}

type RegisterReq struct {
	Account   string `json:"account" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	College   string `json:"college"`
	StudentNo string `json:"student_no"`
	Phone     string `json:"phone"`
}

// @Summary 注册
// @Tags 公共
// @Accept json
// @Produce json
// @Param payload body RegisterReq true "注册参数"
// @Success 200 {object} response.Body
// @Router /public/register [post]
func Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	if req.Gender != "" && !(req.Gender == "male" || req.Gender == "female" || req.Gender == "secret") {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	var exists models.User
	if err := store.DB().Where("account = ?", req.Account).First(&exists).Error; err == nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "账号已存在"))
		return
	}
	hp, err := password.Hash(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "注册失败"))
		return
	}
	var role models.Role
	if err := store.DB().Where("code = ?", "student").First(&role).Error; err != nil {
		role = models.Role{Name: "学生", Code: "student"}
		_ = store.DB().Create(&role).Error
	}
	u := models.User{Account: req.Account, Password: hp, Name: req.Name, Gender: req.Gender, College: req.College, StudentNo: req.StudentNo, Phone: req.Phone, RoleID: role.ID}
	if err := store.DB().Create(&u).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "注册失败"))
		return
	}
	token, err := jwt.Sign(config.Default().JWT.Secret, u.ID, "student", config.Default().JWT.Expires)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "注册失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(map[string]any{"token": token}))
}

// @Summary 我的信息
// @Tags 学生
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /student/me [get]
func MyProfile(c *gin.Context) {
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	var clubCount int64
	var actCount int64
	_ = store.DB().Model(&models.Membership{}).Where("user_id = ? AND status = ?", u.ID, "approved").Count(&clubCount).Error
	_ = store.DB().Model(&models.Attendance{}).Where("user_id = ?", u.ID).Count(&actCount).Error
	res := map[string]any{
		"account":        u.Account,
		"name":           u.Name,
		"gender":         u.Gender,
		"college":        u.College,
		"student_no":     u.StudentNo,
		"phone":          u.Phone,
		"avatar":         u.Avatar,
		"role":           u.Role.Code,
		"club_count":     clubCount,
		"activity_count": actCount,
	}
	c.JSON(http.StatusOK, response.Success(res))
}

type MyUpdateReq struct {
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	College   string `json:"college"`
	StudentNo string `json:"student_no"`
	Phone     string `json:"phone"`
	Avatar    string `json:"avatar"`
}

// @Summary 更新我的信息
// @Tags 学生
// @Accept json
// @Produce json
// @Param payload body MyUpdateReq true "更新参数"
// @Security Bearer
// @Success 200 {object} response.Body
// @Router /student/me [put]
func UpdateMyProfile(c *gin.Context) {
	cu, _ := c.Get("currentUser")
	u := cu.(*models.User)
	var req MyUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	if req.Gender != "" && !(req.Gender == "male" || req.Gender == "female" || req.Gender == "secret") {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	updates := map[string]any{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Gender != "" {
		updates["gender"] = req.Gender
	}
	if req.College != "" {
		updates["college"] = req.College
	}
	if req.StudentNo != "" {
		updates["student_no"] = req.StudentNo
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if len(updates) == 0 {
		c.JSON(http.StatusOK, response.Success(map[string]any{"ok": true}))
		return
	}
	if err := store.DB().Model(&models.User{}).Where("id = ?", u.ID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "更新失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(map[string]any{"ok": true}))
}
