package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/givetrack/givetrack/internal/service"
	"github.com/givetrack/givetrack/internal/util"
)

// ProjectHandler 项目处理器。
type ProjectHandler struct {
	projectSvc *service.ProjectService
}

func NewProjectHandler(projectSvc *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectSvc: projectSvc}
}

// List 项目列表（分类/状态筛选 + 分页）。
func (h *ProjectHandler) List(c *gin.Context) {
	page, ps := util.NormalizePage(atoi(c.Query("page")), atoi(c.Query("page_size")))
	if c.Query("limit") != "" {
		ps = atoi(c.Query("limit"))
	}
	list, total, totalPages, err := h.projectSvc.List(c.Query("category"), c.DefaultQuery("status", "approved"), page, ps)
	if err != nil {
		util.FailError(c, err)
		return
	}
	util.OK(c, gin.H{
		"projects":   list,
		"total":      total,
		"page":       page,
		"limit":      ps,
		"totalPages": totalPages,
	})
}

// GetDetail 项目详情。
func (h *ProjectHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, 40000, "invalid project id")
		return
	}
	project, donations, updates, err := h.projectSvc.GetDetail(uint(id))
	if err != nil {
		util.FailError(c, err)
		return
	}
	// 捐赠列表脱敏
	donationViews := make([]gin.H, 0, len(donations))
	for _, d := range donations {
		donorName := d.User.RealName
		if donorName == "" {
			donorName = d.User.Username
		}
		if d.IsAnonymous {
			donorName = "爱心人士"
		}
		donationViews = append(donationViews, gin.H{
			"id":          d.ID,
			"amount":      d.Amount,
			"isAnonymous": d.IsAnonymous,
			"message":     d.Message,
			"createdAt":   d.CreatedAt,
			"donorName":   donorName,
		})
	}
	util.OK(c, gin.H{"project": project, "donations": donationViews, "updates": updates})
}

// Create 组织发布项目。
func (h *ProjectHandler) Create(c *gin.Context) {
	var req service.CreateProjectInput
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 42200, err.Error())
		return
	}
	p, err := h.projectSvc.Create(c.GetUint("user_id"), req)
	if err != nil {
		util.FailError(c, err)
		return
	}
	util.Created(c, gin.H{"project": p})
}

// MyProjects 组织自己的项目。
func (h *ProjectHandler) MyProjects(c *gin.Context) {
	list, err := h.projectSvc.MyProjects(c.GetUint("user_id"))
	if err != nil {
		util.FailError(c, err)
		return
	}
	util.OK(c, gin.H{"projects": list})
}

// Updates 项目进展列表。
func (h *ProjectHandler) Updates(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, 40000, "invalid project id")
		return
	}
	// 复用详情查询中的进展
	_, _, updates, err := h.projectSvc.GetDetail(uint(id))
	if err != nil {
		util.FailError(c, err)
		return
	}
	util.OK(c, gin.H{"updates": updates})
}

// CreateUpdate 上传执行进展。
func (h *ProjectHandler) CreateUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, 40000, "invalid project id")
		return
	}
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content"`
		Images  string `json:"images"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 42200, err.Error())
		return
	}
	u, err := h.projectSvc.CreateUpdate(c.GetUint("user_id"), uint(id), req.Title, req.Content, req.Images)
	if err != nil {
		util.FailError(c, err)
		return
	}
	util.Created(c, gin.H{"update": u})
}

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
