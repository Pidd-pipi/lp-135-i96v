package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/givetrack/givetrack/internal/constants"
	"github.com/givetrack/givetrack/internal/model"
	"github.com/givetrack/givetrack/internal/repository"
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
	p := &model.Project{
		OrganizationID: org.ID,
		Title:          in.Title,
		Description:    in.Description,
		Category:       in.Category,
		TargetAmount:   in.TargetAmount,
		ExecutionPlan:  in.ExecutionPlan,
		Status:         constants.ProjectPending,
	}
	if in.StartDate != "" {
		if t, err := time.Parse("2006-01-02", in.StartDate); err == nil {
			p.StartDate = &t
		}
	}
	if in.EndDate != "" {
		if t, err := time.Parse("2006-01-02", in.EndDate); err == nil {
			p.EndDate = &t
		}
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

// 项目管理业务错误。
var (
	ErrProjectNotOwner      = errors.New("forbidden: not your project")
	ErrProjectCompleted     = errors.New("project completed and cannot be modified")
	ErrTargetBelowRaised    = errors.New("target amount cannot be lower than current amount")
	ErrInvalidProjectStatus = errors.New("invalid project status for this action")
	ErrInvalidCategory      = errors.New("invalid project category")
	ErrInvalidDateRange     = errors.New("end date cannot be earlier than start date")
	ErrInvalidDate          = errors.New("invalid date, expected format 2006-01-02")
)

// UpdateProjectInput 组织修改项目入参（标题不可改）。
type UpdateProjectInput struct {
	Description   string  `json:"description"`
	Category      string  `json:"category" binding:"required"`
	TargetAmount  float64 `json:"targetAmount" binding:"required,gt=0"`
	ExecutionPlan string  `json:"executionPlan"`
	StartDate     string  `json:"startDate"`
	EndDate       string  `json:"endDate"`
}

// validCategories 允许的项目分类。
var validCategories = map[string]bool{
	constants.CategoryEducation:   true,
	constants.CategoryElderly:     true,
	constants.CategoryMedical:     true,
	constants.CategoryDisaster:    true,
	constants.CategoryEnvironment: true,
	constants.CategoryOther:       true,
}

// UpdateProject 负责人修改自己的项目。
// 已完成项目不可改；目标金额不得低于已筹金额；仅项目所属组织可操作。
func (s *ProjectService) UpdateProject(userID, projectID uint, in UpdateProjectInput) (*model.Project, error) {
	p, org, err := s.loadOwnedProject(userID, projectID)
	if err != nil {
		return nil, err
	}
	start, end, err := validateProjectUpdate(p, in)
	if err != nil {
		return nil, err
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
	s.logger.Info("project updated", "projectId", p.ID, "orgId", org.ID)
	return p, nil
}

// validateProjectUpdate 校验编辑入参与项目当前状态，返回解析后的起止日期。
func validateProjectUpdate(p *model.Project, in UpdateProjectInput) (*time.Time, *time.Time, error) {
	if p.Status == constants.ProjectCompleted {
		return nil, nil, ErrProjectCompleted
	}
	if !validCategories[in.Category] {
		return nil, nil, ErrInvalidCategory
	}
	if in.TargetAmount < p.CurrentAmount {
		return nil, nil, ErrTargetBelowRaised
	}
	return parseDateRange(in.StartDate, in.EndDate)
}

// validateStatusAction 校验状态变更动作与当前状态是否匹配。
func validateStatusAction(currentStatus, action string) (string, error) {
	switch action {
	case "pause":
		if currentStatus != constants.ProjectApproved {
			return "", ErrInvalidProjectStatus
		}
		return constants.ProjectPaused, nil
	case "resume":
		if currentStatus != constants.ProjectPaused {
			return "", ErrInvalidProjectStatus
		}
		return constants.ProjectApproved, nil
	default:
		return "", ErrInvalidProjectStatus
	}
}

// ChangeProjectStatus 筹款中项目停募（approved -> paused）或重新开放（paused -> approved）。
func (s *ProjectService) ChangeProjectStatus(userID, projectID uint, action string) (*model.Project, error) {
	p, org, err := s.loadOwnedProject(userID, projectID)
	if err != nil {
		return nil, err
	}
	newStatus, err := validateStatusAction(p.Status, action)
	if err != nil {
		return nil, err
	}
	p.Status = newStatus
	if err := s.projectRepo.Update(p); err != nil {
		return nil, err
	}
	s.logger.Info("project status changed", "projectId", p.ID, "orgId", org.ID, "action", action, "status", p.Status)
	return p, nil
}

// loadOwnedProject 加载项目并校验归属。
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
		return nil, nil, ErrProjectNotOwner
	}
	return p, org, nil
}

// parseDateRange 解析起止日期，空串表示清空；结束日期不得早于开始日期。
func parseDateRange(startStr, endStr string) (*time.Time, *time.Time, error) {
	var start, end *time.Time
	if startStr != "" {
		t, err := time.Parse("2006-01-02", startStr)
		if err != nil {
			return nil, nil, ErrInvalidDate
		}
		start = &t
	}
	if endStr != "" {
		t, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			return nil, nil, ErrInvalidDate
		}
		end = &t
	}
	if start != nil && end != nil && end.Before(*start) {
		return nil, nil, ErrInvalidDateRange
	}
	return start, end, nil
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
