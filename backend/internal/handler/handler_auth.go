package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/givetrack/givetrack/internal/service"
	"github.com/givetrack/givetrack/internal/util"
)

// AuthHandler 认证处理器。
type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// RegisterRequest 注册请求。
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=72"`
	Role     string `json:"role"`
	RealName string `json:"realName"`
	Phone    string `json:"phone"`
}

// Register 注册。
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 42200, err.Error())
		return
	}
	user, org, token, err := h.authSvc.Register(service.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		RealName: req.RealName,
		Phone:    req.Phone,
	})
	if err != nil {
		util.FailError(c, err)
		return
	}
	util.Created(c, gin.H{"token": token, "user": user, "organization": org})
}

// LoginRequest 登录请求。
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 登录。
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 42200, err.Error())
		return
	}
	user, org, token, err := h.authSvc.Login(req.Username, req.Password)
	if err != nil {
		util.Fail(c, http.StatusUnauthorized, 40100, "invalid username or password")
		return
	}
	util.OK(c, gin.H{"token": token, "user": user, "organization": org})
}

// Me 当前用户信息。
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetUint("user_id")
	util.OK(c, gin.H{"userId": userID, "username": c.GetString("username"), "role": c.GetString("role")})
}

// UpdateProfileRequest 更新资料请求。
type UpdateProfileRequest struct {
	RealName string `json:"realName"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
}

// UpdateProfile 更新个人资料。
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 42200, err.Error())
		return
	}
	user, err := h.authSvc.UpdateProfile(c.GetUint("user_id"), req.RealName, req.Phone, req.Avatar)
	if err != nil {
		util.FailError(c, err)
		return
	}
	util.OK(c, gin.H{"user": user})
}
