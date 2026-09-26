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
