package service

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/givetrack/givetrack/internal/constants"
	"github.com/givetrack/givetrack/internal/model"
	"github.com/givetrack/givetrack/internal/repository"
	"github.com/givetrack/givetrack/internal/util"
)

// ProjectService 项目管理服务。
type ProjectService struct {
	projectRepo  *repository.ProjectRepository
	updateRepo   *repository.ProjectUpdateRepository
	orgRepo      *repository.OrganizationRepository
	donRepo      *repository.DonationRepository
	logger       *slog.Logger
}

func NewProjectService(projectRepo *repository.ProjectRepository, updateRepo *repository.ProjectUpdateRepository, orgRepo *repository.OrganizationRepository, donRepo *repository.DonationRepository, logger *slog.Logger) *ProjectService {
	return &ProjectService{projectRepo: projectRepo, updateRepo: updateRepo, orgRepo: orgRepo, donRepo: donRepo, logger: logger}
}

// ProjectWithProgress 带进度百分比的项目视图。
type ProjectWithProgress struct {
	model.Project
	Progress int `json:"progress"`
}

func withProgress(p *model.Project) ProjectWithProgress {
	progress := 0
	if p.TargetAmount > 0 {
		progress = int(p.CurrentAmount / p.TargetAmount * 100)
		if progress > 100 {
			progress = 100
		}
	}
	return ProjectWithProgress{Project: *p, Progress: progress}
}

// List 项目列表。
func (s *ProjectService) List(category, status string, page, pageSize int) ([]ProjectWithProgress, int64, int, error) {
	list, total, err := s.projectRepo.List(category, status, page, pageSize)
	if err != nil {
		return nil, 0, 0, err
	}
	out := make([]ProjectWithProgress, 0, len(list))
	for i := range list {
		out = append(out, withProgress(&list[i]))
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	return out, total, totalPages, nil
}

// GetDetail 项目详情 + 捐赠记录 + 进展。
func (s *ProjectService) GetDetail(id uint) (ProjectWithProgress, []model.Donation, []model.ProjectUpdate, error) {
	p, err := s.projectRepo.FindByID(id)
	if err != nil {
		return ProjectWithProgress{}, nil, nil, err
	}
	donations, err := s.donRepo.ListByProject(id, 20)
	if err != nil {
		return ProjectWithProgress{}, nil, nil, err
	}
	updates, err := s.updateRepo.ListByProject(id)
	if err != nil {
		return ProjectWithProgress{}, nil, nil, err
	}
	return withProgress(p), donations, updates, nil
}

// Create 组织发布项目。
func (s *ProjectService) Create(userID uint, in CreateProjectInput) (*model.Project, error) {
	org, err := s.orgRepo.FindByUserID(userID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("organization not found")
	}
	if err != nil {
		return nil, err
	}
	if org.Status != constants.OrgApproved {
		return nil, fmt.Errorf("organization not approved")
	}
	if !validCategory(in.Category) {
		return nil, fmt.Errorf("invalid category")
	}
	start, end, err := parseProjectDates(in.StartDate, in.EndDate)
	if err != nil {
		return nil, err
	}
	p := &model.Project{
		OrganizationID: org.ID,
		Title:          in.Title,
		Description:    in.Description,
		Category:       in.Category,
		TargetAmount:   in.TargetAmount,
		ExecutionPlan:  in.ExecutionPlan,
		StartDate:      start,
		EndDate:        end,
		Status:         constants.ProjectPending,
	}
	if err := s.projectRepo.Create(p); err != nil {
		return nil, err
	}
	s.logger.Info("project created", "projectId", p.ID, "orgId", org.ID)
	return p, nil
}

// MyProjects 组织自己的项目。
func (s *ProjectService) MyProjects(userID uint) ([]ProjectWithProgress, error) {
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	list, err := s.projectRepo.ListByOrg(org.ID)
	if err != nil {
		return nil, err
	}
	out := make([]ProjectWithProgress, 0, len(list))
	for i := range list {
		out = append(out, withProgress(&list[i]))
	}
	return out, nil
}

// CreateUpdate 上传项目执行进展。
func (s *ProjectService) CreateUpdate(userID, projectID uint, title, content, images string) (*model.ProjectUpdate, error) {
	p, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, err
	}
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if p.OrganizationID != org.ID {
		return nil, fmt.Errorf("forbidden: not your project")
	}
	u := &model.ProjectUpdate{ProjectID: projectID, Title: title, Content: content, Images: images}
	if err := s.updateRepo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

// UpdateProjectInput 组织修改项目入参（仅允许修改项目介绍等白名单字段）。
type UpdateProjectInput struct {
	Description   string  `json:"description"`
	Category      string  `json:"category" binding:"required"`
	TargetAmount  float64 `json:"targetAmount" binding:"required,gt=0"`
	ExecutionPlan string  `json:"executionPlan"`
	StartDate     string  `json:"startDate"`
	EndDate       string  `json:"endDate"`
}

// 可修改项目信息的状态：待审核、已通过、临时停募。已完成/已驳回不可改。
var editableStatuses = map[string]bool{
	constants.ProjectPending:  true,
	constants.ProjectApproved: true,
	constants.ProjectPaused:   true,
}

// Update 项目负责人修改自己项目的介绍、分类、目标金额、执行计划和起止日期。
// 目标金额不得低于已筹金额；已完成项目不允许修改。
func (s *ProjectService) Update(userID, projectID uint, in UpdateProjectInput) (*model.Project, error) {
	if !validCategory(in.Category) {
		return nil, fmt.Errorf("invalid category")
	}
	start, end, err := parseProjectDates(in.StartDate, in.EndDate)
	if err != nil {
		return nil, err
	}

	p, _, err := s.loadOwnedProject(userID, projectID)
	if err != nil {
		return nil, err
	}
	if !editableStatuses[p.Status] {
		if p.Status == constants.ProjectCompleted {
			return nil, fmt.Errorf("project already completed and cannot be edited")
		}
		return nil, fmt.Errorf("project in status %s cannot be edited", p.Status)
	}
	// 用“分”做整数比较，避免浮点误差导致目标金额略低于已筹金额。
	if cents(in.TargetAmount) < cents(p.CurrentAmount) {
		return nil, fmt.Errorf("target amount cannot be lower than current raised amount")
	}

	p.Description = in.Description
	p.Category = in.Category
	p.TargetAmount = in.TargetAmount
	p.ExecutionPlan = in.ExecutionPlan
	p.StartDate = start
	p.EndDate = end
	if err := s.projectRepo.Update(p); err != nil {
		return nil, err
	}
	s.logger.Info("project updated", "projectId", p.ID, "orgId", p.OrganizationID)
	return p, nil
}

// ChangeStatus 项目负责人停募（approved -> paused）或重新开放（paused -> approved）。
// 已完成项目与其他非法状态流转一律拒绝。
func (s *ProjectService) ChangeStatus(userID, projectID uint, target string) (*model.Project, error) {
	if target != constants.ProjectPaused && target != constants.ProjectApproved {
		return nil, fmt.Errorf("unsupported project status: %s", target)
	}
	p, _, err := s.loadOwnedProject(userID, projectID)
	if err != nil {
		return nil, err
	}
	switch target {
	case constants.ProjectPaused:
		if p.Status != constants.ProjectApproved {
			return nil, fmt.Errorf("only fundraising projects can be paused, current status: %s", p.Status)
		}
	case constants.ProjectApproved:
		if p.Status != constants.ProjectPaused {
			return nil, fmt.Errorf("only paused projects can be reopened, current status: %s", p.Status)
		}
	}
	p.Status = target
	if err := s.projectRepo.Update(p); err != nil {
		return nil, err
	}
	s.logger.Info("project status changed", "projectId", p.ID, "orgId", p.OrganizationID, "status", target)
	return p, nil
}

// loadOwnedProject 加载项目并校验归属；非本人项目返回 ErrForbidden。
func (s *ProjectService) loadOwnedProject(userID, projectID uint) (*model.Project, *model.Organization, error) {
	p, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, nil, err
	}
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil {
		return nil, nil, err
	}
	if p.OrganizationID != org.ID {
		return nil, nil, fmt.Errorf("%w: not your project", util.ErrForbidden)
	}
	return p, org, nil
}

func validCategory(c string) bool {
	switch c {
	case constants.CategoryEducation, constants.CategoryElderly, constants.CategoryMedical,
		constants.CategoryDisaster, constants.CategoryEnvironment, constants.CategoryOther:
		return true
	}
	return false
}

// parseProjectDates 解析起止日期（空串表示不设置），并校验结束日期不早于开始日期。
func parseProjectDates(startStr, endStr string) (*time.Time, *time.Time, error) {
	var start, end *time.Time
	if startStr != "" {
		t, err := time.Parse("2006-01-02", startStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid startDate, expected format YYYY-MM-DD")
		}
		start = &t
	}
	if endStr != "" {
		t, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid endDate, expected format YYYY-MM-DD")
		}
		end = &t
	}
	if start != nil && end != nil && end.Before(*start) {
		return nil, nil, fmt.Errorf("endDate cannot be earlier than startDate")
	}
	return start, end, nil
}

// cents 将金额换算为最小单位“分”的整数，用于精确比较。
func cents(v float64) int64 {
	return int64(math.Round(v * 100))
}

// CreateProjectInput 项目创建入参。
type CreateProjectInput struct {
	Title         string  `json:"title" binding:"required"`
	Description   string  `json:"description"`
	Category      string  `json:"category" binding:"required"`
	TargetAmount  float64 `json:"targetAmount" binding:"required,gt=0"`
	ExecutionPlan string  `json:"executionPlan"`
	StartDate     string  `json:"startDate"`
	EndDate       string  `json:"endDate"`
}
