package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/givetrack/givetrack/internal/service"
	"github.com/givetrack/givetrack/internal/util"
)

// AdminHandler 平台审核处理器。
type AdminHandler struct {
	adminSvc *service.AdminService
}

func NewAdminHandler(adminSvc *service.AdminService) *AdminHandler {
	return &AdminHandler{adminSvc: adminSvc}
}

// PendingProjects 待审核项目。
func (h *AdminHandler) PendingProjects(c *gin.Context) {
	list, err := h.adminSvc.PendingProjects()
	if err != nil {
		util.FailError(c, err)
		return
	}
	util.OK(c, gin.H{"projects": list})
}

// ReviewProject 审核项目。
func (h *AdminHandler) ReviewProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, 40000, "invalid project id")
		return
	}
	var req struct {
		Status  string `json:"status" binding:"required"`
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 42200, err.Error())
		return
	}
	p, err := h.adminSvc.ReviewProject(c.GetUint("user_id"), uint(id), req.Status, req.Comment)
	if err != nil {
		util.FailError(c, err)
		return
	}
	util.OK(c, gin.H{"project": p})
}

// PendingOrganizations 待审核组织。
func (h *AdminHandler) PendingOrganizations(c *gin.Context) {
	list, err := h.adminSvc.PendingOrganizations()
	if err != nil {
		util.FailError(c, err)
		return
	}
	util.OK(c, gin.H{"organizations": list})
}

// ReviewOrganization 审核组织。
func (h *AdminHandler) ReviewOrganization(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, 40000, "invalid organization id")
		return
	}
	var req struct {
		Status  string `json:"status" binding:"required"`
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 42200, err.Error())
		return
	}
	o, err := h.adminSvc.ReviewOrganization(c.GetUint("user_id"), uint(id), req.Status, req.Comment)
	if err != nil {
		util.FailError(c, err)
		return
	}
	util.OK(c, gin.H{"organization": o})
}
